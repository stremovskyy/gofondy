/*
 * MIT License
 *
 * Copyright (c) 2022 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
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

	"github.com/google/uuid"
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/manager"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
	"github.com/stremovskyy/gofondy/utils"
)

type gateway struct {
	manager manager.FondyManager
	options *models.Options
}

func New(options *models.Options) FondyGateway {
	return &gateway{
		manager: manager.NewManager(options),
		options: options,
	}
}

func (g *gateway) VerificationLink(account *models.MerchantAccount, invoiceId uuid.UUID, email *string, note string, code consts.CurrencyCode) (*url.URL, error) {
	fondyVerificationAmount := g.options.VerificationAmount * 100
	lf := strconv.FormatFloat(g.options.VerificationLifeTime.Seconds(), 'f', 2, 64)
	cbu := g.options.CallbackBaseURL + g.options.CallbackUrl

	request := &models.FondyRequestObject{
		MerchantID:        utils.StringRef(account.MerchantID),
		DesignID:          &account.MerchantDesignID,
		Verification:      utils.StringRef("Y"),
		MerchantData:      utils.StringRef(note + "/card verification"),
		Amount:            utils.StringRef(fmt.Sprintf("%d", fondyVerificationAmount)),
		OrderID:           utils.StringRef(invoiceId.String()),
		OrderDesc:         utils.StringRef(g.options.VerificationDescription),
		Lifetime:          utils.StringRef(lf),
		RequiredRectoken:  utils.StringRef("Y"),
		Currency:          utils.StringRef(code.String()),
		ServerCallbackURL: utils.StringRef(cbu),
		SenderEmail:       email,
	}

	raw, err := g.manager.Verify(request, account)
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

func (g *gateway) Status(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID: invoiceRequest.GetMerchantIDString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
	}

	raw, err := g.manager.Status(request, invoiceRequest.Merchant)
	if err != nil {
		return nil, models.NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := models.UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	err = fondyResponse.Error()
	if err != nil {
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *gateway) Refund(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID: invoiceRequest.GetMerchantIDString(),
		Amount:     invoiceRequest.GetAmountString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
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

func (g *gateway) SplitRefund(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error) {
	request := &models_v2.Order{
		MerchantID: invoiceRequest.Merchant.MerchantIDInt(),
		Amount:     invoiceRequest.GetAmountString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
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

func (g *gateway) Split(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error) {
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

	orderData, err := g.Status(invoiceRequest)
	if err != nil {
		return nil, err
	}

	if !orderData.Captured() {
		return nil, errors.New("split accounts problem: order is not captured")
	}

	order := &models_v2.Order{
		MerchantID:  invoiceRequest.Merchant.MerchantIDInt(),
		Amount:      invoiceRequest.GetAmountString(),
		OrderID:     invoiceRequest.GetInvoiceIDString(),
		Currency:    utils.StringRef(string(consts.CurrencyCodeUAH)),
		OrderType:   utils.StringRef("settlement"),
		Rectoken:    invoiceRequest.PaymentCardToken,
		OperationID: invoiceRequest.GetInvoiceIDString(),
		OrderDesc:   invoiceRequest.GetDescriptionString(),
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

func (g *gateway) Payment(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID: invoiceRequest.GetMerchantIDString(),
		Amount:     invoiceRequest.GetAmountString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
		Preauth:    utils.StringRef("N"),
		OrderDesc:  invoiceRequest.GetDescriptionString(),
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

func (g *gateway) Hold(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID: invoiceRequest.GetMerchantIDString(),
		Amount:     invoiceRequest.GetAmountString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
		Preauth:    utils.StringRef("Y"),
		OrderDesc:  invoiceRequest.GetDescriptionString(),
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

func (g *gateway) Capture(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	request := &models.FondyRequestObject{
		MerchantID: invoiceRequest.GetMerchantIDString(),
		Amount:     invoiceRequest.GetAmountString(),
		OrderID:    invoiceRequest.GetInvoiceIDString(),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
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

func (g *gateway) Credit(invoiceRequest *models.InvoiceRequest) (*models.Order, error) {
	captureAmount := invoiceRequest.Amount * 100
	addData := make(map[string]string)
	addData["acc"] = invoiceRequest.InvoiceID.String()

	request := &models.FondyRequestObject{
		MerchantID:       &invoiceRequest.Merchant.MerchantID,
		Amount:           utils.StringRef(fmt.Sprintf("%.f", captureAmount)),
		OrderID:          utils.StringRef(invoiceRequest.InvoiceID.String()),
		Currency:         utils.StringRef(string(consts.CurrencyCodeUAH)),
		ReceiverRectoken: invoiceRequest.WithdrawalCardToken,
		AdditionalData:   addData,
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
		return nil, models.NewAPIError(802, "Fondy Gate Response Failure", err, request, raw)
	}

	return &fondyResponse.Response, nil
}
