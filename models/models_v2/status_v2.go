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
	Order Order `json:"order"`
}

type Transaction struct {
	Status                        string      `json:"status"`
	OdbRef                        interface{} `json:"odb_ref"`
	DocNo                         interface{} `json:"doc_no"`
	Currency                      string      `json:"currency"`
	SettlementResponseCode        string      `json:"settlement_response_code"`
	MerchantID                    int64       `json:"merchant_id"`
	ID                            int64       `json:"id"`
	SettlementCurrency            string      `json:"settlement_currency"`
	SettlementFee                 float64     `json:"settlement_fee"`
	Fee                           int64       `json:"fee"`
	ReversalAmount                int64       `json:"reversal_amount"`
	SettlementAmount              float64     `json:"settlement_amount"`
	SettlementStatus              string      `json:"settlement_status"`
	Amount                        int64       `json:"amount"`
	SettlementResponseDescription string      `json:"settlement_response_description"`
	ParentTranID                  interface{} `json:"parent_tran_id"`
	Receiver                      Receiver    `json:"receiver"`
	Payouttime                    string      `json:"payouttime"`
	Type                          string      `json:"type"`
}
