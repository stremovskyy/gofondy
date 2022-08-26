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

package models_v2

import "encoding/json"

func UnmarshalStatus(data []byte) (Status, error) {
	var r Status
	err := json.Unmarshal(data, &r)
	return r, err
}

type Status struct {
	Order Order `json:"order,omitempty"`
}

type Transaction struct {
	Status                        string      `json:"status,omitempty"`
	OdbRef                        interface{} `json:"odb_ref,omitempty"`
	DocNo                         interface{} `json:"doc_no,omitempty"`
	Currency                      string      `json:"currency,omitempty"`
	SettlementResponseCode        string      `json:"settlement_response_code,omitempty"`
	MerchantID                    int64       `json:"merchant_id,omitempty"`
	ID                            int64       `json:"id,omitempty"`
	SettlementCurrency            string      `json:"settlement_currency,omitempty"`
	SettlementFee                 float64     `json:"settlement_fee,omitempty"`
	Fee                           int64       `json:"fee,omitempty"`
	ReversalAmount                int64       `json:"reversal_amount,omitempty"`
	SettlementAmount              float64     `json:"settlement_amount,omitempty"`
	SettlementStatus              string      `json:"settlement_status,omitempty"`
	Amount                        int64       `json:"amount,omitempty"`
	SettlementResponseDescription string      `json:"settlement_response_description,omitempty"`
	ParentTranID                  interface{} `json:"parent_tran_id,omitempty"`
	Receiver                      Receiver    `json:"receiver,omitempty"`
	Payouttime                    string      `json:"payouttime,omitempty"`
	Type                          string      `json:"type,omitempty"`
}
