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

package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy"
	"github.com/karmadon/gofondy/examples"
	"github.com/karmadon/gofondy/models"
)

func main() {
	fondyGateway := gofondy.New(models.DefaultOptions())

	merchAccount := &models.MerchantAccount{
		MerchantID:     examples.MerchantId,
		MerchantKey:    examples.MerchantKey,
		MerchantString: "Test Merchant",
	}

	invoiceId := uuid.MustParse("4eeefd24-d9d0-48a2-9f21-31138c0f6447")

	captureAmount := float64(1)

	capturePayment, err := fondyGateway.CapturePayment(merchAccount, &invoiceId, &captureAmount)
	if err != nil {
		log.Fatal(err)
	}

	if *capturePayment.ResponseStatus == "success" {
		fmt.Printf("Order (%s) Capture Status: %s\n", capturePayment.OrderID.String(), *capturePayment.CaptureStatus)
	} else {
		fmt.Printf("Error: %s\n", capturePayment.ErrorMessage)
	}
}
