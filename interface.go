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

package gofondy

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
)

type FondyGateway interface {
	VerificationLink(account *models.MerchantAccount, invoiceId uuid.UUID, email *string, note string, code consts.CurrencyCode) (*url.URL, error)
	Status(account *models.MerchantAccount, invoiceId *uuid.UUID) (*models.Order, error)
	Payment(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.Order, error)
	Hold(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64, token string) (*models.Order, error)
	Capture(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models.Order, error)
	Refund(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models.Order, error)
	Split(account *models.MerchantAccount, invoiceId *uuid.UUID, token string) (*models_v2.Order, error)
	SplitRefund(account *models.MerchantAccount, invoiceId *uuid.UUID, amount *float64) (*models_v2.Order, error)
}
