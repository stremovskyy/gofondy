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

package manager

import (
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
	"github.com/stremovskyy/gofondy/utils"
)

type FondyManager interface {
	StraightPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	MobileHoldPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	MobileStraightPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	CapturePayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	Verify(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	Status(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	HoldPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	Withdraw(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	RefundPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	SplitRefund(order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error)
	SplitPayment(order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error)
}

type manager struct {
	client  Client
	options *models.Options
}

func NewManager(options *models.Options) *manager {
	return &manager{
		options: options,
		client: NewClient(&ClientOptions{
			Timeout:         options.Timeout,
			KeepAlive:       options.KeepAlive,
			MaxIdleConns:    options.MaxIdleConns,
			IdleConnTimeout: options.IdleConnTimeout,
		}),
	}
}

func NewManagerWithClient(options *models.Options, client Client) *manager {
	return &manager{
		options: options,
		client:  client,
	}
}

func (m *manager) HoldPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	request.MerchantData = utils.StringRef("hold/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.client.payment(consts.FondyURLRecurring, request, merchantAccount)
}

func (m *manager) StraightPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	request.MerchantData = utils.StringRef("straight/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())

	return m.client.payment(consts.FondyURLRecurring, request, merchantAccount)
}

func (m *manager) MobileHoldPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	request.Preauth = utils.StringRef("Y")
	request.MerchantData = utils.StringRef("mobile/hold/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())
	request.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	request.MerchantID = &merchantAccount.MerchantID

	return m.client.payment(consts.Fondy3DSecureS1, request, merchantAccount)
}

func (m *manager) MobileStraightPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	request.Preauth = utils.StringRef("N")
	request.MerchantData = utils.StringRef("mobile/straight/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())
	request.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	request.MerchantID = &merchantAccount.MerchantID

	return m.client.payment(consts.Fondy3DSecureS1, request, merchantAccount)
}

func (m *manager) Withdraw(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	request.MerchantData = utils.StringRef("withdraw/" + merchantAccount.MerchantAddedDescription + request.AdditionalDataString())
	request.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	request.MerchantID = &merchantAccount.MerchantID

	return m.client.withdraw(consts.FondyURLP2PCredit, request, merchantAccount)
}

func (m *manager) CapturePayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.payment(consts.FondyURLCapture, request, merchantAccount)
}

func (m *manager) RefundPayment(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.payment(consts.FondyURLRefund, request, merchantAccount)
}

func (m *manager) Status(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.payment(consts.FondyURLStatus, request, merchantAccount)
}

func (m *manager) Verify(request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.payment(consts.FondyURLGetVerification, request, merchantAccount)
}

func (m *manager) SplitPayment(order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.split(consts.FondySettlement, order, merchantAccount)
}

func (m *manager) SplitRefund(order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.client.split(consts.FondyURLRefund, order, merchantAccount)
}
