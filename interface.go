/*
 * Project: banker
 * File: interface.go (4/29/23, 4:37 PM)
 *
 * Copyright (C) Megakit Systems 2017-2023, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package gofondy

import (
	"net/url"

	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
)

type FondyGateway interface {
	V1() V1
	V2() V2
	ID() ID
}

type V1 interface {
	VerificationLink(invoiceRequest *models.InvoiceRequest) (*url.URL, error)
	Status(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
	Payment(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
	Hold(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
	Capture(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
	Refund(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
	Credit(invoiceRequest *models.InvoiceRequest) (*models.Order, error)
}

type V2 interface {
	SplitRefund(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error)
	Split(invoiceRequest *models.InvoiceRequest) (*models_v2.Order, error)
}

type ID interface {
	Status(*models.IDStatusRequest) (*models.FondyClientStatusResponse, error)
	Limits(*models.IDStatusRequest) (*models.FondyBalance, error)
}
