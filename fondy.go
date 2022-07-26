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
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy/consts"
	"github.com/karmadon/gofondy/manager"
	"github.com/karmadon/gofondy/models"
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

	if fondyResponse.Response.ResponseStatus != nil && *fondyResponse.Response.ResponseStatus != "success" {
		return nil, models.NewAPIError(802, *fondyResponse.Response.ErrorMessage, err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *gateway) Refund(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *int) (*models.OrderData, error) {
	refundAmount := *amount * 100

	request := &models.RequestObject{
		Amount:   utils.StringRef(fmt.Sprintf("%d", refundAmount)),
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

	if fondyResponse.Response.ResponseStatus != nil && *fondyResponse.Response.ResponseStatus != "success" {
		return nil, models.NewAPIError(802, *fondyResponse.Response.ErrorMessage, err, request, raw)
	}

	return &fondyResponse.Response, nil
}
