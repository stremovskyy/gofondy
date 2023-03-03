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

package models_v2

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/stremovskyy/gofondy/consts"
)

func UnmarshalResponse(data []byte) (ResponseWrapper, error) {
	var r ResponseWrapper
	err := json.Unmarshal(data, &r)
	return r, err
}

type ResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	Version   string `json:"version"`
	Data      []byte `json:"data"`
	Signature string `json:"signature"`
}

func (w *ResponseWrapper) SignIsValid(key string) bool {
	if w == nil {
		return false
	}

	s := key + "|" + string(w.Response.Data)
	h := sha1.New()
	h.Write([]byte(s))
	calculated := fmt.Sprintf("%x", h.Sum(nil))

	return calculated == w.Response.Signature
}

func (w *ResponseWrapper) Error() error {
	order, err := w.Order()
	if err != nil {
		return fmt.Errorf("error while getting order: %w", err)
	}

	if order.ResponseStatus != consts.FondyResponseStatusSuccess && order.ResponseStatus != consts.FondyResponseStatusCreated {
		return fmt.Errorf("order status is %s, code: %s", order.ResponseStatus, order.ResponseCode)
	}

	return nil
}

func (w *ResponseWrapper) Order() (*Order, error) {
	var wrapper OrderWrapper

	err := json.Unmarshal(w.Response.Data, &wrapper)
	if err != nil {
		return nil, err
	}

	return &wrapper.Order, nil
}
