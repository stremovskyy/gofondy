/*
 * MIT License
 *
 * Copyright (c) 2024 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
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
	"crypto/sha1"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/stremovskyy/gofondy/utils"
)

type FondyClientStatusRequest struct {
	MerchantID       *string `json:"merchant_id"`
	IPN              *string `json:"ipn,omitempty"`
	InternalPassport *string `json:"internal_passport,omitempty"`
	IDCard           *string `json:"id_card,omitempty"`

	Signature *string `json:"signature"`
}

// Sign - adds signature for request using provided key
func (r *FondyClientStatusRequest) Sign(key string, isDebug bool) error {
	if r.Signature != nil {
		r.Signature = nil
	}

	s := key + "|"

	values := reflect.ValueOf(*r)
	types := values.Type()
	preFiltered := map[string]string{}

	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Signature" {
			continue
		}

		t := values.Field(i).Interface()

		if t != nil {
			s, ok := t.(*string)
			if ok && s != nil {
				preFiltered[types.Field(i).Name] = *s
			}
		}
	}

	keys := make([]string, 0, len(preFiltered))
	for k := range preFiltered {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	final := make([]string, 0, len(preFiltered))
	for _, k := range keys {
		if value, exists := preFiltered[k]; exists && value != "" {
			final = append(final, value)
		}
	}

	s += strings.Join(final, "|")

	h := sha1.New()
	h.Write([]byte(s))
	r.Signature = utils.StringRef(fmt.Sprintf("%x", h.Sum(nil)))

	if isDebug {
		fmt.Println("[GO FONDY ID] Sign String: ", s)
		fmt.Println("[GO FONDY ID] Signature: ", *r.Signature)
	}

	return nil
}
