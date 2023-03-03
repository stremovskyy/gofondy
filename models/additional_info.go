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

package models

import (
	"encoding/json"

	"github.com/stremovskyy/gofondy/consts"
)

func UnmarshalAdditionalInfo(data []byte) (AdditionalInfo, error) {
	var r AdditionalInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

type AdditionalInfo struct {
	CaptureStatus           consts.FondyCaptureStatus `json:"capture_status,omitempty"`
	CaptureAmount           *float64                  `json:"capture_amount,omitempty"`
	ReservationData         interface{}               `json:"reservation_data"`
	TransactionID           *int64                    `json:"transaction_id,omitempty"`
	BankResponseCode        *string                   `json:"bank_response_code"`
	BankResponseDescription *string                   `json:"bank_response_description"`
	ClientFee               *int64                    `json:"client_fee,omitempty"`
	SettlementFee           *float64                  `json:"settlement_fee,omitempty"`
	BankName                *string                   `json:"bank_name,omitempty"`
	BankCountry             *string                   `json:"bank_country,omitempty"`
	CardType                *string                   `json:"card_type,omitempty"`
	CardProduct             *string                   `json:"card_product,omitempty"`
	CardCategory            *string                   `json:"card_category,omitempty"`
	Timeend                 *string                   `json:"timeend,omitempty"`
	IpaddressV4             *string                   `json:"ipaddress_v4,omitempty"`
	PaymentMethod           *string                   `json:"payment_method,omitempty"`
}
