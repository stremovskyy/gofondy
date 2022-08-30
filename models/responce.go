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

package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/karmadon/gofondy/consts"
)

func UnmarshalFondyResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

type Response struct {
	Response ResponseObject `json:"response"`
}

func (r *Response) Error() error {
	if r == nil {
		return errors.New("response object is nil")
	}

	if r.Response.ResponseStatus != consts.FondyResponseStatusSuccess {
		errString := "Fondy Response is not successful"

		if r.Response.ErrorMessage != nil {
			errString += " Message: " + *r.Response.ErrorMessage
		}

		if r.Response.ErrorCode != nil {
			errString += " Code: " + strconv.Itoa(int(*r.Response.ErrorCode))
		}

		return errors.New(errString)
	}

	return nil
}

type ResponseObject struct {
	Target         string                     `json:"target"`
	ResponseURL    *string                    `json:"response_url"`
	ResponseStatus consts.FondyResponseStatus `json:"response_status"`
	Pending        bool                       `json:"pending"`
	OrderData      Order                      `json:"order_data"`
	APIVersion     string                     `json:"api_version"`
	PaymentID      *string                    `json:"payment_id"`
	CheckoutURL    *string                    `json:"checkout_url"`
	ErrorMessage   *string                    `json:"error_message"`
	ErrorCode      *int64                     `json:"error_code"`
	RequestID      *string                    `json:"request_id"`
}
