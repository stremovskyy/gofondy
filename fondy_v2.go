/*
 * MIT License
 *
 * Copyright (c) 2024 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gofondy

import (
	"errors"
	"fmt"

	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/manager"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
	"github.com/stremovskyy/gofondy/utils"
)

var v2 V2

type fondyV2 struct {
	manager manager.FondyManager
	options *models.Options
}

func (g *gateway) V2() V2 {
	if v2 == nil {
		v2 = &fondyV2{
			manager: g.manager,
			options: g.options,
		}
	}

	return v2
}

func (g *fondyV2) SplitRefund(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error) {
	request := &models_v2.Order{
		MerchantID:        invoiceRequest.Merchant.MerchantIDInt(),
		Amount:            invoiceRequest.GetAmountString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	raw, err := g.manager.SplitRefund(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := models_v2.UnmarshalResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	order, err := fondyResponse.Order()
	if err != nil {
		return nil, err
	}

	if order.ReverseStatus != consts.FondyReverseStatusSuccess && order.ReverseStatus != consts.FondyReverseStatusApproved {
		err = fmt.Errorf("reverse status is %s, (%s)", order.ReverseStatus, *order.ResponseDescription)
		return nil, models.NewAPIError(803, "Fondy Gate Response Failure", err, request, raw)
	}

	return order, nil
}

func (g *fondyV2) Split(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error) {
	err := invoiceRequest.Merchant.SplitAccounts.Error()
	if err != nil {
		return nil, errors.New("split accounts problem " + err.Error())
	}

	if !invoiceRequest.Merchant.IsTechnical {
		return nil, errors.New("split accounts problem: only technical accounts can split")
	}

	if len(invoiceRequest.Merchant.SplitAccounts) == 0 {
		return nil, errors.New("split accounts problem: no split accounts")
	}

	request := &models.FondyRequestObject{
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	rawStatus, err := g.manager.Status(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed", err, request, rawStatus)
	}

	fondyStatusResponse, err := models.UnmarshalStatusResponse(*rawStatus)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, request, rawStatus)
	}

	if !fondyStatusResponse.Response.Captured() {
		return nil, errors.New("split accounts problem: order is not captured")
	}

	order := &models_v2.Order{
		MerchantID:        invoiceRequest.Merchant.MerchantIDInt(),
		Amount:            invoiceRequest.GetAmountString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		OrderType:         utils.StringRef("settlement"),
		Rectoken:          invoiceRequest.PaymentCardToken,
		OperationID:       invoiceRequest.GetInvoiceIDString(),
		OrderDesc:         invoiceRequest.GetDescriptionString(),
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	raw, err := g.manager.SplitPayment(order, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http splitRequest failed", err, nil, raw)
	}

	fondyResponse, err := models_v2.UnmarshalResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, nil, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, nil, raw)
	}

	return fondyResponse.Order()
}
