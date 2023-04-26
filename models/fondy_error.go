/*
 * Project: banker
 * File: fondy_error.go (4/26/23, 1:18 PM)
 *
 * Copyright (C) Megakit Systems 2017-2023, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package models

import "strconv"

type FondyError struct {
	ErrorCode    int64
	ErrorMessage string
	IsFatal      bool
}

func NewFatalFondyError(errorCode int, errorMessage string) *FondyError {
	return &FondyError{ErrorCode: int64(errorCode), ErrorMessage: errorMessage, IsFatal: true}
}

func NewFondyError(errorCode string, errorMessage string) *FondyError {
	code, _ := strconv.ParseInt(errorCode, 10, 32)
	return &FondyError{ErrorCode: code, ErrorMessage: errorMessage, IsFatal: false}
}

func (e *FondyError) Error() string {
	return "[FONDY] " + e.ErrorMessage + " (" + strconv.Itoa(int(e.ErrorCode)) + ")"
}
