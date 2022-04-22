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
	"time"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy"
)

func main() {
	options := &gofondy.Options{
		Timeout:                 30 * time.Second,
		KeepAlive:               30 * time.Second,
		MaxIdleConns:            10,
		IdleConnTimeout:         20 * time.Second,
		VerificationAmount:      1,
		VerificationDescription: "Verification Test",
		VerificationLifeTime:    600 * time.Second,
		CallbackUrl:             FondyVerificationCallbackURL,
		DesignId:                DesignId,
		MerchantId:              MerchantId,
		MerchantKey:             MerchantKey,
	}

	fondyGateway := gofondy.New(options)

	verificationLink, err := fondyGateway.VerificationLink(uuid.New(), nil, "test", gofondy.CurrencyCodeUAH)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Verification link: %s", *verificationLink)
}
