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

package models

import (
	"encoding/json"
	"errors"
	"github.com/stremovskyy/gofondy/consts"
	"strconv"
)

func UnmarshalStatusResponse(data []byte) (StatusResponse, error) {
	var r StatusResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

type StatusResponse struct {
	Response Order `json:"response"`
}

func (r *StatusResponse) Error() error {
	if r == nil {
		return errors.New("response object is nil")
	}

	if r.Response.ResponseStatus != nil && *r.Response.ResponseStatus != consts.FondyResponseStatusSuccess {
		errString := "Fondy Response is not successful"

		if r.Response.ErrorMessage != nil {
			errString += " Message: " + *r.Response.ErrorMessage
		}

		if r.Response.ErrorCode != nil {
			return NewFatalFondyError(int(*r.Response.ErrorCode), errString)
		}

		return NewFatalFondyError(-1, errString)
	}

	if r.Response.ResponseCode != nil {
		if r.Response.ResponseDescription != nil {
			if code, ok := r.Response.ResponseCode.(string); ok {
				if code == "" {
					return nil
				}

				return NewFondyError(code, *r.Response.ResponseDescription)
			}

			if code, ok := r.Response.ResponseCode.(int64); ok {
				return NewFondyError(strconv.FormatInt(code, 10), *r.Response.ResponseDescription)
			}

			if code, ok := r.Response.ResponseCode.(float64); ok {
				return NewFondyError(strconv.FormatFloat(code, 'f', -1, 64), *r.Response.ResponseDescription)
			}

			return NewFondyError("-1", *r.Response.ResponseDescription)
		}
	}

	return nil
}
