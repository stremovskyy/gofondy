/*
 * MIT License
 *
 * Copyright (c) 2022 Anton (karmadon) Stremovskyy <stremovskyy@gmail.com>
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
	"strconv"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy/consts"
	"github.com/karmadon/gofondy/manager"
	"github.com/karmadon/gofondy/models"
	"github.com/karmadon/gofondy/models/models_v2"
	"github.com/karmadon/gofondy/utils"
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

func (g *gateway) VerificationLink(account *models.MerchantAccount, invoiceId uuid.UUID, email *string, note string, code consts.CurrencyCode) (*string, error) {
	fondyVerificationAmount := g.options.VerificationAmount * 100
	lf := strconv.FormatFloat(g.options.VerificationLifeTime.Seconds(), 'f', 2, 64)
	cbu := g.options.CallbackBaseURL + g.options.CallbackUrl

	request := &models.RequestObject{
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

	return fondyResponse.Response.CheckoutURL, nil
}

func (g *gateway) Status(account *models.MerchantAccount, invoiceId *uuid.UUID) (*models.OrderData, error) {
	request := &models.RequestObject{
		MerchantID: &account.MerchantID,
		OrderID:    utils.StringRef(invoiceId.String()),
	}

	raw, err := g.manager.Status(request, account)
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

func (g *gateway) Refund(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models.OrderData, error) {
	refundAmount := *amount * 100

	request := &models.RequestObject{
		MerchantID: &account.MerchantID,
		Amount:     utils.StringRef(fmt.Sprintf("%.f", refundAmount)),
		OrderID:    utils.StringRef(invoiceId.String()),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
	}

	raw, err := g.manager.RefundPayment(request, account)
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

func (g *gateway) SplitRefund(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models_v2.Order, error) {
	refundAmount := *amount * 100

	request := &models_v2.Order{
		MerchantID: account.MerchantIDInt(),
		Amount:     utils.StringRef(fmt.Sprintf("%.f", refundAmount)),
		OrderID:    utils.StringRef(invoiceId.String()),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
	}

	raw, err := g.manager.SplitRefund(request, account)
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

	if order.ReverseStatus == nil || (*order.ReverseStatus != "success" && *order.ReverseStatus != "approved") {
		err = fmt.Errorf("reverse status is %s, (%s)", *order.ReverseStatus, *order.ResponseDescription)
		return nil, models.NewAPIError(803, "Fondy Gate Response Failure", err, request, raw)
	}

	return order, nil
}

func (g *gateway) Split(account *models.MerchantAccount, invoiceId *uuid.UUID, token string) (*models_v2.Order, error) {
	err := account.SplitAccounts.Error()
	if err != nil {
		return nil, errors.New("split accounts problem " + err.Error())
	}

	if !account.IsTechnical {
		return nil, errors.New("split accounts problem: only technical accounts can split")
	}

	if len(account.SplitAccounts) == 0 {
		return nil, errors.New("split accounts problem: no split accounts")
	}

	orderData, err := g.Status(account, invoiceId)
	if err != nil {
		return nil, err
	}

	if !orderData.Captured() {
		return nil, errors.New("split accounts problem: order is not captured")
	}

	amount := orderData.CapturedAmount() * 100

	order := &models_v2.Order{
		MerchantID:  account.MerchantIDInt(),
		Amount:      utils.StringRef(fmt.Sprintf("%.f", amount)),
		OrderID:     utils.StringRef(uuid.NewString()),
		Currency:    utils.StringRef(string(consts.CurrencyCodeUAH)),
		OrderType:   utils.StringRef("settlement"),
		Rectoken:    utils.StringRef(token),
		OperationID: utils.StringRef(invoiceId.String()),
		OrderDesc:   utils.StringRef(account.MerchantString),
	}

	raw, err := g.manager.SplitPayment(order, account)
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

func (g *gateway) Payment(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.OrderData, error) {
	paymentAmount := *amount * 100

	request := &models.RequestObject{
		MerchantID: &account.MerchantID,
		OrderDesc:  utils.StringRef(account.MerchantString),
		Amount:     utils.StringRef(fmt.Sprintf("%.f", paymentAmount)),
		OrderID:    utils.StringRef(invoiceId.String()),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
		Rectoken:   utils.StringRef(token),
		Preauth:    utils.StringRef("N"),
	}

	raw, err := g.manager.StraightPayment(request, account)
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

func (g *gateway) Hold(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.OrderData, error) {
	holdAmount := *amount * 100

	request := &models.RequestObject{
		MerchantID: &account.MerchantID,
		Amount:     utils.StringRef(fmt.Sprintf("%.f", holdAmount)),
		OrderID:    utils.StringRef(invoiceId.String()),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
		Rectoken:   utils.StringRef(token),
		Preauth:    utils.StringRef("Y"),
		OrderDesc:  utils.StringRef(account.MerchantString),
	}

	raw, err := g.manager.HoldPayment(request, account)
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

func (g *gateway) Capture(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models.OrderData, error) {
	captureAmount := *amount * 100

	request := &models.RequestObject{
		MerchantID: &account.MerchantID,
		Amount:     utils.StringRef(fmt.Sprintf("%.f", captureAmount)),
		OrderID:    utils.StringRef(invoiceId.String()),
		Currency:   utils.StringRef(string(consts.CurrencyCodeUAH)),
	}

	raw, err := g.manager.CapturePayment(request, account)
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
