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
	"net/url"
	"strconv"

	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/manager"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/utils"
)

var v1 V1

type fondyV1 struct {
	manager manager.FondyManager
	options *models.Options
}

func (g *gateway) V1() V1 {
	if v1 == nil {
		v1 = &fondyV1{
			manager: g.manager,
			options: g.options,
		}
	}

	return v1
}

func (g *fondyV1) VerificationLink(invoiceRequest *models.InvoiceRequest) (*url.URL, error) {
	fondyVerificationAmount := g.options.VerificationAmount * 100
	lf := strconv.FormatFloat(g.options.VerificationLifeTime.Seconds(), 'f', 2, 64)

	request := &models.FondyRequestObject{
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		DesignID:          &invoiceRequest.Merchant.MerchantDesignID,
		Verification:      utils.StringRef("Y"),
		MerchantData:      utils.StringRef("/card verification"),
		Amount:            utils.StringRef(fmt.Sprintf("%d", fondyVerificationAmount)),
		OrderDesc:         utils.StringRef(g.options.VerificationDescription),
		Lifetime:          utils.StringRef(lf),
		RequiredRectoken:  utils.StringRef("Y"),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	if invoiceRequest.PaymentLifetime != nil {
		sec := int64(invoiceRequest.PaymentLifetime.Seconds())
		request.Lifetime = utils.StringRef(fmt.Sprintf("%d", sec))
	}

	raw, err := g.manager.Verify(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalFondyResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, err
	}

	if fondyResponse.Response.CheckoutURL == nil {
		return nil, models.NewAPIError(802, "No Url In Response", err, request, raw)
	}

	return url.Parse(*fondyResponse.Response.CheckoutURL)
}

func (g *fondyV1) Status(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	raw, err := g.manager.Status(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *fondyV1) Refund(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		Amount:            invoiceRequest.GetAmountString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	raw, err := g.manager.RefundPayment(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "REFUND: API ERROR", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "REFUND: Unmarshal refund response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "REFUND: fondy gate returned an error", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *fondyV1) Payment(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		Amount:            invoiceRequest.GetAmountString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		Preauth:           utils.StringRef("N"),
		OrderDesc:         invoiceRequest.GetDescriptionString(),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	if invoiceRequest.PaymentLifetime != nil {
		sec := int64(invoiceRequest.PaymentLifetime.Seconds())
		request.Lifetime = utils.StringRef(fmt.Sprintf("%d", sec))
	}

	var raw *[]byte
	var err error

	if invoiceRequest.IsMobile() {
		request.RequiredRectoken = utils.StringRef("Y")
		request.Container = invoiceRequest.Container
		raw, err = g.manager.MobileStraightPayment(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	} else {
		if invoiceRequest.PaymentCardToken == nil {
			return nil, errors.New("token is required for web hold")
		}

		request.Rectoken = utils.StringRef(*invoiceRequest.PaymentCardToken)
		raw, err = g.manager.StraightPayment(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	}
	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal hold payment response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *fondyV1) Hold(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:        invoiceRequest.GetMerchantIDString(),
		Amount:            invoiceRequest.GetAmountString(),
		OrderID:           invoiceRequest.GetInvoiceIDString(),
		Currency:          utils.StringRef(string(consts.CurrencyCodeUAH)),
		Preauth:           utils.StringRef("Y"),
		OrderDesc:         invoiceRequest.GetDescriptionString(),
		AdditionalData:    invoiceRequest.AdditionalData,
		ServerCallbackURL: invoiceRequest.ServerCallbackURL,
	}

	if invoiceRequest.PaymentLifetime != nil {
		sec := int64(invoiceRequest.PaymentLifetime.Seconds())
		request.Lifetime = utils.StringRef(fmt.Sprintf("%d", sec))
	}

	var raw *[]byte
	var err error

	if invoiceRequest.IsMobile() {
		request.RequiredRectoken = utils.StringRef("Y")
		request.Container = invoiceRequest.Container
		raw, err = g.manager.MobileHoldPayment(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	} else {
		if invoiceRequest.PaymentCardToken == nil {
			return nil, errors.New("token is required for web hold")
		}

		request.Rectoken = utils.StringRef(*invoiceRequest.PaymentCardToken)
		raw, err = g.manager.HoldPayment(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	}

	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed while holding payment", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal hold payment response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *fondyV1) Capture(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:     invoiceRequest.GetMerchantIDString(),
		Amount:         invoiceRequest.GetAmountString(),
		OrderID:        invoiceRequest.GetInvoiceIDString(),
		Currency:       utils.StringRef(string(consts.CurrencyCodeUAH)),
		AdditionalData: invoiceRequest.AdditionalData,
	}

	raw, err := g.manager.CapturePayment(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed while capturing payment", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal capture response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *fondyV1) Credit(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID:         &invoiceRequest.Merchant.MerchantID,
		Amount:             invoiceRequest.GetAmountString(),
		OrderID:            invoiceRequest.GetInvoiceIDString(),
		Currency:           utils.StringRef(string(consts.CurrencyCodeUAH)),
		ReceiverRectoken:   invoiceRequest.WithdrawalCardToken,
		ReceiverCardNumber: invoiceRequest.WithdrawalCardNumber,
		AdditionalData:     invoiceRequest.AdditionalData,
		ServerCallbackURL:  invoiceRequest.ServerCallbackURL,
	}

	raw, err := g.manager.Withdraw(request, invoiceRequest.Merchant, invoiceRequest.ReservationData)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed while capturing payment", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal capture response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		if errors.As(err, &models.FondyError{}) {
			errorCode := err.(*models.FondyError).ErrorCode
			return &fondyResponse.Response, models.NewAPIError(int(errorCode), "Fondy Gate Response Failure", err, request, raw)
		}
		return &fondyResponse.Response, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}
