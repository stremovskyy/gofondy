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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type gateway struct {
	client  *http.Client
	options *Options
}

func New(options *Options) FondyGateway {
	g := &gateway{options: options}

	dialer := &net.Dialer{
		Timeout:   options.Timeout,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}}

	g.client = &http.Client{Transport: tr}

	return g
}

func (g *gateway) VerificationLink(invoiceId uuid.UUID, email *string, note string, code CurrencyCode) (*string, error) {
	fondyVerificationAmount := g.options.VerificationAmount * 100
	lf := strconv.FormatFloat(g.options.VerificationLifeTime.Seconds(), 'f', 2, 64)
	cbu := g.options.CallbackBaseURL + g.options.CallbackUrl

	request := RequestObject{
		MerchantData:      StringRef(note + "/card verification"),
		Amount:            StringRef(strconv.Itoa(fondyVerificationAmount)),
		OrderID:           StringRef(invoiceId.String()),
		OrderDesc:         StringRef(g.options.VerificationDescription),
		Lifetime:          StringRef(lf),
		Verification:      StringRef("Y"),
		DesignID:          StringRef(g.options.DesignId),
		MerchantID:        StringRef(g.options.MerchantId),
		RequiredRectoken:  StringRef("Y"),
		Currency:          StringRef(code.String()),
		ServerCallbackURL: StringRef(cbu),
		SenderEmail:       email,
	}

	raw, err := g.makeFondyRequest(request, FondyURLGetVerification, false)
	if err != nil {
		return nil, NewAPIError(800, "Http request failed", err, &request, raw)
	}

	fondyResponse, err := UnmarshalFondyResponse(*raw)
	if err != nil {
		return nil, NewAPIError(801, "Unmarshal response fail", err, &request, raw)
	}

	if fondyResponse.Response.CheckoutURL == nil {
		return nil, NewAPIError(802, "No Url In Response", err, &request, raw)
	}

	return fondyResponse.Response.CheckoutURL, nil
}

func (g *gateway) makeFondyRequest(request RequestObject, url FondyURL, credit bool) (*[]byte, error) {
	methodPost := "POST"
	err := request.CreateSignature(g.options.MerchantKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create signature: %w", err)
	}

	jsonValue, err := json.Marshal(NewFondyRequest(request))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	req, err := http.NewRequest(methodPost, url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header = http.Header{
		"User-Agent":   {"Utax driveapp Service/" + Version},
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"X-Request-ID": {uuid.New().String()},
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot copy response buffer: %w", err)
	}
	defer resp.Body.Close()

	return &raw, nil
}
