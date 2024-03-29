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

	"github.com/stremovskyy/gofondy"
	"github.com/stremovskyy/gofondy/examples"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/recorder/redis_recorder"
)

func main() {
	responseRecorder := redis_recorder.NewRedisRecorder(
		redis_recorder.NewDefaultOptions(
			"localhost:6379",
			"",
			12,
		),
	)

	fondyGateway := gofondy.NewWithRecorder(models.DebugDefaultOptions(), responseRecorder)

	merchAccount := &models.MerchantAccount{
		MerchantID:     examples.MerchantId,
		MerchantKey:    examples.MerchantKey,
		MerchantString: "Test Merchant",
	}

	verificationRequest := &models.IDStatusRequest{
		Merchant: merchAccount,
		ID:       "2562515514",
		IDType:   models.IDTypeTIN,
	}

	statusResponse, err := fondyGateway.ID().Status(verificationRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nstatusResponse: %s\n", statusResponse)
	fmt.Printf("\nlimit: %s\n", statusResponse.LimitTill().String())
}
