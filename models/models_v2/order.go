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

type Order struct {
	PaymentID           int64         `json:"payment_id"`
	Fee                 string        `json:"fee"`
	OrderType           string        `json:"order_type"`
	ReversalAmount      string        `json:"reversal_amount"`
	OrderID             string        `json:"order_id"`
	SettlementAmount    string        `json:"settlement_amount"`
	MerchantData        string        `json:"merchant_data"`
	SettlementDate      string        `json:"settlement_date"`
	Transaction         []Transaction `json:"transaction"`
	OperationID         string        `json:"operation_id"`
	OrderStatus         string        `json:"order_status"`
	ResponseDescription string        `json:"response_description"`
	MerchantID          int64         `json:"merchant_id"`
	OrderTime           string        `json:"order_time"`
	ResponseCode        string        `json:"response_code"`
	SettlementCurrency  string        `json:"settlement_currency"`
	ServerCallbackURL   string        `json:"server_callback_url"`
	Rectoken            string        `json:"rectoken"`
	Currency            string        `json:"currency"`
	Amount              string        `json:"amount"`
	ResponseURL         string        `json:"response_url"`
	OrderDesc           string        `json:"order_desc"`
	Receiver            []Receiver    `json:"receiver"`
}
