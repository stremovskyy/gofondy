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

package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/stremovskyy/gofondy"
	"github.com/stremovskyy/gofondy/examples"
	"github.com/stremovskyy/gofondy/models"
)

func main() {
	fondyGateway := gofondy.New(models.DefaultOptions())

	splitA := &models.MerchantAccount{
		MerchantID:      examples.SplitAMerchantId,
		MerchantKey:     examples.SplitAMerchantKey,
		MerchantString:  "Split A Merchant",
		SplitPercentage: 30.0,
	}

	splitB := &models.MerchantAccount{
		MerchantID:      examples.SplitBMerchantId,
		MerchantKey:     examples.SplitBMerchantKey,
		MerchantString:  "Split B Merchant",
		SplitPercentage: 70.0,
	}

	accounts := models.MerchantAccounts{}
	accounts.Add(splitA)
	accounts.Add(splitB)
	err := accounts.Error()
	if err != nil {
		log.Fatal(err)
	}

	techAccount := &models.MerchantAccount{
		MerchantID:     examples.TechMerchantId,
		MerchantKey:    examples.TechMerchantKey,
		MerchantString: "Tech Merchant",
		SplitAccounts:  accounts,
		IsTechnical:    true,
	}

	invoiceId := uuid.MustParse("767f44ef-2997-4623-961f-9ee081ef730f")

	intermediateResponse, err := fondyGateway.V2().Split(techAccount, &invoiceId, examples.CardToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("intermediateResponse: %+v\n", intermediateResponse)
}
