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

package consts

type FondyURL string

const (
	FondyURLGetVerification FondyURL = "https://api.fondy.eu/api/checkout/url/"
	FondyURLStatus          FondyURL = "https://api.fondy.eu/api/status/order_id/"
	FondyURLRecurring       FondyURL = "https://api.fondy.eu/api/recurring/"
	FondyURLP2PCredit       FondyURL = "https://api.fondy.eu/api/p2pcredit/"
	FondyURLRefund          FondyURL = "https://api.fondy.eu/api/reverse/order_id/"
	FondyURLCapture         FondyURL = "https://api.fondy.eu/api/capture/order_id/"
	Fondy3DSecureS1         FondyURL = "https://pay.fondy.eu/api/3dsecure_step1/"
	FondySettlement         FondyURL = "https://pay.fondy.eu/api/settlement"
)

func (t FondyURL) String() string {
	return string(t)
}
