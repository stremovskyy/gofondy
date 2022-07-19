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
)

type gateway struct {
	manager fondyManager
	options *Options
}

func New(options *Options) FondyGateway {
	return &gateway{
		manager: NewManager(options),
		options: options,
	}
}

func (g *gateway) VerificationLink(account *MerchantAccount, invoiceId uuid.UUID, email *string, note string, code CurrencyCode) (*string, error) {
	fondyVerificationAmount := g.options.VerificationAmount * 100
	lf := strconv.FormatFloat(g.options.VerificationLifeTime.Seconds(), 'f', 2, 64)
	cbu := g.options.CallbackBaseURL + g.options.CallbackUrl

	request := &RequestObject{
		MerchantData:      StringRef(note + "/card verification"),
		Amount:            StringRef(fmt.Sprintf("%d", fondyVerificationAmount)),
		OrderID:           StringRef(invoiceId.String()),
		OrderDesc:         StringRef(g.options.VerificationDescription),
		Lifetime:          StringRef(lf),
		RequiredRectoken:  StringRef("Y"),
		Currency:          StringRef(code.String()),
		ServerCallbackURL: StringRef(cbu),
		SenderEmail:       email,
	}

	raw, err := g.manager.Verify(request, account)
	if err != nil {
		return nil, NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := UnmarshalFondyResponse(*raw)
	if err != nil {
		return nil, NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	if fondyResponse.Response.CheckoutURL == nil {
		return nil, NewAPIError(802, "No Url In Response", err, request, raw)
	}

	return fondyResponse.Response.CheckoutURL, nil
}

func (g *gateway) Status(account *MerchantAccount, invoiceId *uuid.UUID) (*OrderData, error) {
	request := &RequestObject{
		OrderID: StringRef(invoiceId.String()),
	}

	raw, err := g.manager.Status(request, account)
	if err != nil {
		return nil, NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	if fondyResponse.Response.ResponseStatus != nil && *fondyResponse.Response.ResponseStatus != "success" {
		return nil, NewAPIError(802, *fondyResponse.Response.ErrorMessage, err, request, raw)
	}

	return &fondyResponse.Response, nil
}

func (g *gateway) Refund(account *MerchantAccount, invoiceId *uuid.UUID, amount *int) (*OrderData, error) {
	refundAmount := *amount * 100

	request := &RequestObject{
		Amount:   StringRef(fmt.Sprintf("%.0f", refundAmount)),
		OrderID:  StringRef(invoiceId.String()),
		Currency: StringRef(string(CurrencyCodeUAH)),
	}

	raw, err := g.manager.RefundPayment(request, account)
	if err != nil {
		return nil, NewAPIError(800, "Http request failed", err, request, raw)
	}

	fondyResponse, err := UnmarshalStatusResponse(*raw)
	if err != nil {
		return nil, NewAPIError(801, "Unmarshal response fail", err, request, raw)
	}

	if fondyResponse.Response.ResponseStatus != nil && *fondyResponse.Response.ResponseStatus != "success" {
		return nil, NewAPIError(802, *fondyResponse.Response.ErrorMessage, err, request, raw)
	}

	return &fondyResponse.Response, nil
}
