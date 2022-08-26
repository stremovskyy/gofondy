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

package models

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

type MerchantGate string

const (
	MerchantGateFondy MerchantGate = "fondy"
)

type MerchantPaymentType string

const (
	MerchantPaymentTypeDriver    MerchantPaymentType = "driver"
	MerchantPaymentTypePassenger MerchantPaymentType = "passenger"
)

type MerchantFlowType string

const (
	MerchantFlowTypeWithdraw MerchantFlowType = "withdraw"
	MerchantFlowTypePayment  MerchantFlowType = "payment"
)

type MerchantAccount struct {
	UUID                     uuid.UUID           `json:"uuid"`
	Name                     string              `json:"name"`
	MerchantString           string              `json:"merchant_string"`
	MerchantAddedDescription string              `json:"merchant_added_description"`
	MerchantPaymentType      MerchantPaymentType `json:"merchant_payment_type"`
	MerchantFlowType         MerchantFlowType    `json:"merchant_flow_type"`
	MerchantGate             MerchantGate        `json:"merchant_gate"`
	MerchantID               string              `json:"merchant_id"`
	MerchantKey              string              `json:"merchant_key"`
	MerchantCreditKey        string              `json:"merchant_credit_key"`
	MerchantDesignID         string              `json:"merchant_design_id"`
	IsTechnical              bool                `json:"is_technical"`
	SplitAccounts            MerchantAccounts    `json:"split_accounts"`
	SplitPercentage          float64             `json:"split_percentage"`
}

func NewMerchantAccount(merchantID string, merchantKey string, merchantCreditKey string) *MerchantAccount {
	return &MerchantAccount{MerchantID: merchantID, MerchantKey: merchantKey, MerchantCreditKey: merchantCreditKey}
}

func (a *MerchantAccount) MerchantIDInt() int64 {
	parseInt, err := strconv.ParseInt(a.MerchantID, 10, 64)
	if err != nil {
		return 0
	}

	return parseInt
}

type MerchantAccounts []MerchantAccount

func (a *MerchantAccounts) Error() error {
	if a == nil {
		return errors.New("empty merchant accounts")
	}

	splitPercent := 0.0

	for _, account := range *a {
		splitPercent += account.SplitPercentage
	}

	if splitPercent != 100.0 {
		return errors.New("split percent sum not equal 100")
	}

	return nil
}

func (a *MerchantAccounts) Add(account *MerchantAccount) {
	*a = append(*a, *account)
}
