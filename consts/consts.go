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

const Version = "0.1.0"

const (
	FondyTimeFormat = "02.01.2006 15:04:05"
)

type FondyTransactionType string

const (
	FondyTransactionTypePurchase     FondyTransactionType = "purchase"
	FondyTransactionTypeReverse      FondyTransactionType = "reverse"
	FondyTransactionTypeVerification FondyTransactionType = "verification"
	FondyTransactionTypeP2PCredit    FondyTransactionType = "p2p credit"
)

func (t FondyTransactionType) String() string {
	return string(t)
}

type FondyResponseStatus string

const (
	FondyResponseStatusSuccess FondyResponseStatus = "success"
	FondyResponseStatusFailure FondyResponseStatus = "failure"
)

func (s FondyResponseStatus) String() string {
	return string(s)
}

type CurrencyCode string

const (
	CurrencyCodeUAH CurrencyCode = "UAH"
)

func (c CurrencyCode) String() string {
	return string(c)
}

type FondyCardType string

const (
	FondyCardTypeVISA       FondyCardType = "VISA"
	FondyCardTypeMASTERCARD FondyCardType = "MASTERCARD"
)

func (t FondyCardType) String() string {
	return string(t)
}

type FondyCaptureStatus string

const (
	FondyCaptureStatusHold     FondyCaptureStatus = "hold"
	FondyCaptureStatusCaptured FondyCaptureStatus = "captured"
)

func (s *FondyCaptureStatus) String() string {
	if s != nil {
		return string(*s)
	}
	return ""
}

// FondyReverseStatus Reversal processing status
type FondyReverseStatus string

const (
	// FondyReverseStatusCreated reversal has been created, but not processed yet
	FondyReverseStatusCreated FondyReverseStatus = "created"
	// FondyReverseStatusDeclined reversal is declined by FONDY payment gateway or by bank or by external payment system
	FondyReverseStatusDeclined FondyReverseStatus = "declined"
	// FondyReverseStatusApproved reversal completed successfully
	FondyReverseStatusApproved FondyReverseStatus = "approved"
)

type Status string

const (
	StatusReversed   Status = "reversed"
	StatusApproved   Status = "approved"
	StatusProcessing Status = "processing"
	StatusDeclined   Status = "declined"
	StatusExpired    Status = "expired"
	StatusCreated    Status = "created"
	StatusCanceled   Status = "canceled"
	StatusCaptured   Status = "captured"
)
