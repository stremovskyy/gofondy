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

package models_v2

import (
	"github.com/stremovskyy/gofondy/consts"
)

type OrderWrapper struct {
	Order Order `json:"order"`
}

type Order struct {
	PaymentID           *int64                     `json:"payment_id,omitempty"`
	Fee                 *string                    `json:"fee,omitempty"`
	OrderType           *string                    `json:"order_type,omitempty"`
	ReversalAmount      *string                    `json:"reversal_amount,omitempty"`
	OrderID             *string                    `json:"order_id,omitempty"`
	SettlementAmount    *string                    `json:"settlement_amount,omitempty"`
	MerchantData        *string                    `json:"merchant_data,omitempty"`
	SettlementDate      *string                    `json:"settlement_date,omitempty"`
	Transaction         []Transaction              `json:"transaction,omitempty"`
	OperationID         *string                    `json:"operation_id,omitempty"`
	OrderStatus         *string                    `json:"order_status,omitempty"`
	ResponseDescription *string                    `json:"response_description,omitempty"`
	MerchantID          int64                      `json:"merchant_id,omitempty"`
	OrderTime           *string                    `json:"order_time,omitempty"`
	ResponseCode        interface{}                `json:"response_code,omitempty"`
	SettlementCurrency  *string                    `json:"settlement_currency,omitempty"`
	ServerCallbackURL   *string                    `json:"server_callback_url,omitempty"`
	Rectoken            *string                    `json:"rectoken,omitempty"`
	Currency            *string                    `json:"currency,omitempty"`
	Amount              *string                    `json:"amount,omitempty"`
	ResponseURL         *string                    `json:"response_url,omitempty"`
	OrderDesc           *string                    `json:"order_desc,omitempty"`
	Receiver            []Receiver                 `json:"receiver,omitempty"`
	ReverseStatus       consts.FondyReverseStatus  `json:"reverse_status,omitempty"`
	ResponseStatus      consts.FondyResponseStatus `json:"response_status,omitempty"`
	ReverseID           *string                    `json:"reverse_id,omitempty"`
	TransactionID       *string                    `json:"transaction_id,omitempty"`
}

func (o *Order) AddReceiver(receiver *Receiver) {
	o.Receiver = append(o.Receiver, *receiver)
}
