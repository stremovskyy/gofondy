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
	"fmt"
	"strconv"

	"github.com/karmadon/gofondy/consts"
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
			errString += " Code: " + strconv.Itoa(int(*r.Response.ErrorCode))
		}

		return errors.New(errString)
	}

	if r.Response.ResponseCode != nil {
		if r.Response.ResponseDescription != nil {
			if code, ok := r.Response.ResponseCode.(string); ok {
				if code == "" {
					return nil
				}

				return errors.New(fmt.Sprintf("code: %s: %s", code, *r.Response.ResponseDescription))
			}

			if code, ok := r.Response.ResponseCode.(int64); ok {
				return errors.New(fmt.Sprintf("code: %d: %s", code, *r.Response.ResponseDescription))
			}

			return fmt.Errorf("~ %#v: %s", r.Response.ResponseCode, *r.Response.ResponseDescription)
		}
	}

	return nil
}
