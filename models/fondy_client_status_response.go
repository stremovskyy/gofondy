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
	"fmt"
	"time"
)

type FondyClientStatusResponse struct {
	IsIdentified bool          `json:"is_identified"`
	IPN          *string       `json:"ipn"`
	Balance      *FondyBalance `json:"balance,omitempty"`
	Error        *string       `json:"error,omitempty"`
}

type FondyBalance struct {
	CurrentLimit float64 `json:"current_limit"`
	UsedLimit    float64 `json:"used_limit"`
	CurrentDate  string  `json:"current_date"`
}

func (r *FondyClientStatusResponse) IsError() bool {
	return r.Error != nil
}

func (r *FondyClientStatusResponse) GetError() error {
	if r.Error != nil {
		return fmt.Errorf("api client fondy error: %s", *r.Error)
	}

	return nil
}

func (r *FondyClientStatusResponse) IsInFondyDB() bool {
	return r.IsIdentified && r.IPN != nil
}

func (r *FondyClientStatusResponse) LeftLimit() float64 {
	return r.Balance.CurrentLimit - r.Balance.UsedLimit
}

func (r *FondyClientStatusResponse) LimitTill() *time.Time {
	t, err := time.Parse("2006-01", r.Balance.CurrentDate)
	if err != nil {
		return nil
	}

	t = t.AddDate(0, 1, 0)

	return &t
}
