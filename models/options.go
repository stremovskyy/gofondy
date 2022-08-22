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
	"time"
)

type Options struct {
	Timeout                 time.Duration
	KeepAlive               time.Duration
	MaxIdleConns            int
	IdleConnTimeout         time.Duration
	VerificationAmount      int
	VerificationDescription string
	VerificationLifeTime    time.Duration

	CallbackBaseURL string
	CallbackUrl     string
}

func DefaultOptions() *Options {
	return &Options{
		Timeout:                 time.Second * 10,
		KeepAlive:               time.Second * 10,
		MaxIdleConns:            10,
		IdleConnTimeout:         time.Second * 10,
		VerificationAmount:      1,
		VerificationDescription: "Verification Test",
		VerificationLifeTime:    600 * time.Second,

		CallbackBaseURL: "http://localhost:8080",
		CallbackUrl:     "/callback",
	}
}