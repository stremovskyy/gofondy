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

package models_v2

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type KeepingZeroFloat float64

func (f KeepingZeroFloat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.1f", float64(f))), nil
}

type Request struct {
	Version   KeepingZeroFloat `json:"version"` // Idiotic version parsing in Fondy API
	Data      string           `json:"data"`
	Signature string           `json:"signature"`
}

func NewRequest(order *Order) *RequestWrapper {
	if order == nil {
		return nil
	}

	dataEncoded, _ := encodeToBase64(OrderWrapper{Order: *order})

	return &RequestWrapper{
		Request{Version: 2.0, Data: dataEncoded},
	}
}

type RequestWrapper struct {
	Request Request `json:"request"`
}

func (w *RequestWrapper) Sign(key string) *RequestWrapper {
	if w == nil {
		return nil
	}

	s := key + "|" + w.Request.Data
	h := sha1.New()
	h.Write([]byte(s))
	w.Request.Signature = fmt.Sprintf("%x", h.Sum(nil))

	return w
}

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}

	err = encoder.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
