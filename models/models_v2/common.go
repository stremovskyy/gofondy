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

type Receiver struct {
	Requisites Requisites `json:"requisites"`
	Type       string     `json:"type"`
}

func NewMerchantReceiver(requisites *Requisites) *Receiver {
	return &Receiver{Requisites: *requisites, Type: "merchant"}
}

type Requisites struct {
	Amount                int64   `json:"amount"`
	SettlementDescription *string `json:"settlement_description,omitempty"`
	MerchantID            *string `json:"merchant_id,omitempty"` // TODO: fondy couldn't decide string or int64
	Account               *int64  `json:"account,omitempty"`
	Okpo                  *int64  `json:"okpo,omitempty"`
	JurName               *string `json:"jur_name,omitempty"`
	Rectoken              *string `json:"rectoken,omitempty"`
	CardNumber            *int64  `json:"card_number,omitempty"`
}

func NewMerchantRequisites(amount int64, merchantID *string, settlementDescription *string) *Requisites {
	return &Requisites{Amount: amount, SettlementDescription: settlementDescription, MerchantID: merchantID}
}
