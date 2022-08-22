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
		OrderID: utils.StringRef(invoiceId.String()),
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
		Amount:   utils.StringRef(fmt.Sprintf("%.f", refundAmount)),
		OrderID:  utils.StringRef(invoiceId.String()),
		Currency: utils.StringRef(string(consts.CurrencyCodeUAH)),
	}

	raw, err := g.manager.RefundPayment(request, account)
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

func (g *gateway) Split(account *models.MerchantAccount, invoiceId *uuid.UUID, cardToken string) (*models_v2.Response, error) {
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

	order := models_v2.Order{
		Amount:      orderData.Amount,
		OrderID:     utils.StringRef(invoiceId.String()),
		Currency:    utils.StringRef(string(consts.CurrencyCodeUAH)),
		OrderType:   utils.StringRef("settlement"),
		Rectoken:    utils.StringRef(cardToken),
		OperationID: utils.StringRef(invoiceId.String()),
		Receiver:    []models_v2.Receiver{},
	}

	wholeAmount, err := strconv.ParseFloat(*orderData.Amount, 64)
	if err != nil {
		return nil, errors.New("split accounts problem: amount parse error")
	}

	splitAmountSum := 0.0

	for _, merchantAccount := range account.SplitAccounts {
		splitAmount := wholeAmount * merchantAccount.SplitPercentage / 100
		merchantReceiver := models_v2.NewMerchantReceiver(models_v2.NewMerchantRequisites(int64(splitAmount), &merchantAccount.MerchantID, &merchantAccount.MerchantAddedDescription))
		order.Receiver = append(order.Receiver, *merchantReceiver)
		splitAmountSum += splitAmount
	}

	if splitAmountSum != wholeAmount {
		return nil, fmt.Errorf("order %s split accounts problem: split amount sum %f != whole amount %f", orderData.OrderID.String(), splitAmountSum, wholeAmount)
	}

	splitRequest := &models_v2.SplitRequest{Order: order}

	raw, err := g.manager.SplitPayment(splitRequest, account)
	if err != nil {
		return nil, models.NewAPIError(800, "Http splitRequest failed", err, nil, raw)
	}

	fondyResponse, err := models_v2.UnmarshalResponse(*raw)
	if err != nil {
		return nil, models.NewAPIError(801, "Unmarshal response fail", err, nil, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *gateway) PaymentByToken(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.OrderData, error) {
	paymentAmount := *amount * 100

	request := &models.RequestObject{
		Amount:   utils.StringRef(fmt.Sprintf("%.f", paymentAmount)),
		OrderID:  utils.StringRef(invoiceId.String()),
		Currency: utils.StringRef(string(consts.CurrencyCodeUAH)),
		Rectoken: utils.StringRef(token),
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

func (g *gateway) HoldPaymentByToken(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.OrderData, error) {
	holdAmount := *amount * 100

	request := &models.RequestObject{
		Amount:   utils.StringRef(fmt.Sprintf("%.f", holdAmount)),
		OrderID:  utils.StringRef(invoiceId.String()),
		Currency: utils.StringRef(string(consts.CurrencyCodeUAH)),
		Rectoken: utils.StringRef(token),
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

func (g *gateway) CapturePayment(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models.OrderData, error) {
	captureAmount := *amount * 100

	request := &models.RequestObject{
		Amount:   utils.StringRef(fmt.Sprintf("%.f", captureAmount)),
		OrderID:  utils.StringRef(invoiceId.String()),
		Currency: utils.StringRef(string(consts.CurrencyCodeUAH)),
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
