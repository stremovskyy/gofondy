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
	"crypto/sha1"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy/consts"
)

type OrderData struct {
	ErrorMessage            *string                      `json:"error_message"`
	Rrn                     *string                      `json:"rrn"`
	MaskedCard              *string                      `json:"masked_card"`
	SenderCellPhone         *string                      `json:"sender_cell_phone"`
	ResponseSignatureString *string                      `json:"response_signature_string"`
	ResponseStatus          *consts.FondyResponseStatus  `json:"response_status"`
	SenderAccount           *string                      `json:"sender_account"`
	Fee                     *string                      `json:"fee"`
	RectokenLifetime        *string                      `json:"rectoken_lifetime"`
	ReversalAmount          *string                      `json:"reversal_amount"`
	CaptureStatus           *consts.FondyCaptureStatus   `json:"capture_status"`
	SettlementAmount        *string                      `json:"settlement_amount"`
	ActualAmount            *string                      `json:"actual_amount"`
	OrderStatus             *consts.Status               `json:"order_status"`
	ResponseDescription     *string                      `json:"response_description"`
	VerificationStatus      *string                      `json:"verification_status"`
	OrderTime               *string                      `json:"order_time"`
	ActualCurrency          *consts.CurrencyCode         `json:"actual_currency"`
	OrderID                 *uuid.UUID                   `json:"order_id"`
	ParentOrderID           *string                      `json:"parent_order_id"`
	MerchantData            *string                      `json:"merchant_data"`
	TranType                *consts.FondyTransactionType `json:"tran_type"`
	Eci                     *string                      `json:"eci"`
	SettlementDate          *string                      `json:"settlement_date"`
	PaymentSystem           *string                      `json:"payment_system"`
	Rectoken                *string                      `json:"rectoken"`
	ApprovalCode            *string                      `json:"approval_code"`
	MerchantID              *int                         `json:"merchant_id"`
	SettlementCurrency      *consts.CurrencyCode         `json:"settlement_currency"`
	PaymentID               *int                         `json:"payment_id"`
	ProductID               *string                      `json:"product_id"`
	Currency                *consts.CurrencyCode         `json:"currency"`
	CardBin                 interface{}                  `json:"card_bin"`
	ResponseCode            interface{}                  `json:"response_code"`
	CardType                *consts.FondyCardType        `json:"card_type"`
	Amount                  *string                      `json:"amount"`
	SenderEmail             *string                      `json:"sender_email"`
	Signature               *string                      `json:"signature"`
	ErrorCode               *int64                       `json:"error_code"`
	FeeOplata               *string                      `json:"fee_oplata"`
}

func (o *OrderData) SignValid(merchantKey string) bool {
	if o.Signature == nil {
		return false
	}
	s := merchantKey + "|"

	values := reflect.ValueOf(*o)
	types := values.Type()
	preFiltered := map[string]string{}

	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Signature" || types.Field(i).Name == "ResponseSignatureString" {
			continue
		}
		t := values.Field(i).Interface()
		if t != nil {
			s, ok := t.(*string)
			if ok && s != nil && len(*s) > 0 {
				preFiltered[types.Field(i).Name] = *s
			} else if str, ok := t.(fmt.Stringer); ok && len(str.String()) > 0 {
				preFiltered[types.Field(i).Name] = str.String()
			} else if num, ok := t.(float64); ok {
				preFiltered[types.Field(i).Name] = fmt.Sprintf("%.0f", num)
			} else if dig, ok := t.(*int); ok {
				preFiltered[types.Field(i).Name] = strconv.Itoa(*dig)
			}
		}
	}

	keys := make([]string, 0, len(preFiltered))
	for k := range preFiltered {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	final := make([]string, 0, len(preFiltered))
	for _, k := range keys {
		final = append(final, preFiltered[k])
	}

	sk := strings.Join(final, "|")
	s += sk

	h := sha1.New()
	h.Write([]byte(s))

	createdSignature := fmt.Sprintf("%x", h.Sum(nil))

	if createdSignature != *o.Signature {
		return false
	}

	return true
}

func (o *OrderData) Captured() bool {
	if o == nil {
		return false
	}

	if o.OrderStatus != nil && *o.OrderStatus == consts.StatusApproved && (o.ReversalAmount == nil || *o.ReversalAmount == "0") {
		return o.FeeOplata != nil && *o.FeeOplata != "0"
	}

	return false
}

func (o *OrderData) Reversed() bool {
	if o == nil {
		return false
	}

	if o.OrderStatus != nil && *o.OrderStatus == consts.StatusReversed && o.ReversalAmount != nil && *o.ReversalAmount != "" {
		return true
	}

	if o.ReversalAmount != nil && *o.ReversalAmount != "" && *o.ReversalAmount != "0" {
		return true
	}

	return false
}

func (o *OrderData) Undefined() bool {
	if o == nil && o.OrderStatus == nil {
		return false
	}

	if *o.OrderStatus == consts.StatusApproved && o.ReversalAmount != nil && *o.ReversalAmount != "" {
		return true
	}

	return false
}

func (o *OrderData) UncompletedHold() bool {
	if o == nil && o.OrderStatus == nil {
		return false
	}

	if *o.OrderStatus == consts.StatusApproved && (o.ReversalAmount == nil || *o.ReversalAmount == "") {
		return o.FeeOplata == nil || *o.FeeOplata == "0"
	}

	return false
}

func (o *OrderData) Declined() bool {
	if o == nil && o.OrderStatus == nil {
		return false
	}

	return *o.OrderStatus == consts.StatusDeclined
}

func (o *OrderData) Expired() bool {
	if o == nil && o.OrderStatus == nil {
		return false
	}

	return *o.OrderStatus == consts.StatusExpired
}

func (o *OrderData) RealAmount() float64 {
	if o == nil || o.Amount == nil {
		return 0
	}

	amount, err := strconv.ParseFloat(*o.Amount, 64)
	if err != nil {
		return 0
	}

	return amount / 100
}

func (o *OrderData) ReversedAmount() float64 {
	if o == nil || o.ReversalAmount == nil {
		return 0
	}

	amount, err := strconv.ParseFloat(*o.ReversalAmount, 64)
	if err != nil {
		return 0
	}

	return amount / 100
}

func (o *OrderData) SplitedAmount() float64 {
	if o == nil || o.SettlementAmount == nil {
		return 0
	}

	amount, err := strconv.ParseFloat(*o.SettlementAmount, 64)
	if err != nil {
		return 0
	}

	return amount / 100
}
