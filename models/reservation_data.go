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
	"github.com/stremovskyy/gofondy/utils"
	"strconv"
)

type ReservationData struct {
	Phonemobile       *string `json:"phonemobile,omitempty"`
	Account           *string `json:"account,omitempty"`
	Uuid              *string `json:"uuid,omitempty"`
	ReceiverInn       *string `json:"receiver_inn,omitempty"`
	ReceiverPan       *string `json:"receiver_pan,omitempty"`
	ReceiverToken     *string `json:"receiver_token,omitempty"`
	PurchasePaymentId *string `json:"purchase_payment_id,omitempty"`
}

func NewReservationDataForPaymentID(purchasePaymentId *int64) *ReservationData {
	if purchasePaymentId == nil {
		return nil
	}

	id := strconv.Itoa(int(*purchasePaymentId))
	return &ReservationData{PurchasePaymentId: &id}
}

func NewReservationDataForReceiverToken(receiverToken *string) *ReservationData {
	return &ReservationData{ReceiverToken: receiverToken}
}

func (r *ReservationData) Base64Encoded() *string {
	b, _ := utils.Base64StructEncode(r)
	return &b
}
