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

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type fondyManager interface {
	StraightPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	MobileHoldPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	MobileStraightPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	CapturePayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	Verify(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	Status(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	HoldPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	Withdraw(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
	RefundPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error)
}

type manager struct {
	client  *http.Client
	options *Options
}

func NewManager(options *Options) *manager {
	m := &manager{options: options}

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

	m.client = &http.Client{Transport: tr}

	return m
}

func (m *manager) HoldPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.Preauth = StringRef("Y")
	request.MerchantData = StringRef("hold/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.payment(FondyURLRecurring, request, merchantAccount)
}

func (m *manager) StraightPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.Preauth = StringRef("N")
	request.MerchantData = StringRef("straight/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.payment(FondyURLRecurring, request, merchantAccount)
}

func (m *manager) MobileHoldPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.Preauth = StringRef("Y")
	request.MerchantData = StringRef("mobile/hold/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.payment(Fondy3DSecureS1, request, merchantAccount)
}

func (m *manager) MobileStraightPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.Preauth = StringRef("N")
	request.MerchantData = StringRef("mobile/straight/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.payment(Fondy3DSecureS1, request, merchantAccount)
}

func (m *manager) CapturePayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.final(FondyURLCapture, request, merchantAccount)
}

func (m *manager) RefundPayment(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.final(FondyURLRefund, request, merchantAccount)
}

func (m *manager) Withdraw(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.MerchantData = StringRef("withdraw/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.withdraw(FondyURLP2PCredit, request, merchantAccount)
}

func (m *manager) Verify(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	request.DesignID = &merchantAccount.MerchantDesignID
	request.Verification = StringRef("Y")

	return m.payment(FondyURLGetVerification, request, merchantAccount)
}

func (m *manager) Status(request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.final(FondyURLStatus, request, merchantAccount)
}

func (m *manager) payment(url FondyURL, request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, true)
}

func (m *manager) info(url FondyURL, request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, false)
}

func (m *manager) withdraw(url FondyURL, request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.do(url, request, true, merchantAccount, true)
}

func (m *manager) final(url FondyURL, request *RequestObject, merchantAccount *MerchantAccount) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, false)
}

func (m *manager) do(url FondyURL, request *RequestObject, credit bool, merchantAccount *MerchantAccount, addOrderDescription bool) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	request.MerchantID = &merchantAccount.MerchantID

	if addOrderDescription {
		request.OrderDesc = StringRef(merchantAccount.MerchantString)
	}

	if credit {
		err := request.Sign(merchantAccount.MerchantCreditKey)
		if err != nil {
			return nil, errors.Errorf("cannot sign request with credit key: %v", err)
		}
	} else {
		err := request.Sign(merchantAccount.MerchantKey)
		if err != nil {
			return nil, errors.Errorf("cannot sign request with merchant key: %v", err)
		}
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
		"User-Agent":   {"GOFONDY/" + Version},
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"X-Request-ID": {requestID},
	}

	resp, err := m.client.Do(req)
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
