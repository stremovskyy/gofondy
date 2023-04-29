/*
 * Project: banker
 * File: invoice_request.go (4/29/23, 9:11 AM)
 *
 * Copyright (C) Megakit Systems 2017-2023, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package models

import (
	"fmt"
	"github.com/google/uuid"
)

type InvoiceRequest struct {
	InvoiceID           uuid.UUID
	Merchant            *MerchantAccount
	Amount              float64
	PaymentCardToken    *string
	WithdrawalCardToken *string
	ReservationData     *ReservationData
	Container           *string
}

func (i *InvoiceRequest) GetInvoiceIDString() *string {
	if i.InvoiceID == uuid.Nil {
		return nil
	}

	id := i.InvoiceID.String()
	return &id
}

func (i *InvoiceRequest) GetAmountString() *string {
	amount := fmt.Sprintf("%d", int64(i.Amount*100))
	return &amount
}

func (i *InvoiceRequest) GetMerchantIDString() *string {
	if i.Merchant == nil {
		return nil
	}

	return &i.Merchant.MerchantID
}

func (i *InvoiceRequest) GetDescriptionString() *string {
	if i.Merchant == nil {
		return nil
	}

	return &i.Merchant.MerchantAddedDescription
}

func (i *InvoiceRequest) IsMobile() bool {
	return i.Container != nil
}
