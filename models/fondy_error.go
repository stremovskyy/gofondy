/*
 * Project: banker
 * File: fondy_error.go (5/22/23, 12:22 PM)
 *
 * Copyright (C) Megakit Systems 2017-2023, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package models

import (
	"github.com/stremovskyy/gofondy/fondy_status"
	"strconv"
)

type FondyError struct {
	ErrorCode    fondy_status.StatusCode
	ErrorMessage string
	IsFatal      bool
}

func NewFatalFondyError(errorCode int, errorMessage string) *FondyError {
	return &FondyError{ErrorCode: fondy_status.StatusCode(errorCode), ErrorMessage: errorMessage, IsFatal: true}
}

func NewFondyError(statusCode fondy_status.StatusCode, errorMessage string) *FondyError {
	return &FondyError{ErrorCode: statusCode, ErrorMessage: errorMessage, IsFatal: false}
}

func (e FondyError) Error() string {
	return "[FONDY] " + e.ErrorMessage + " (" + strconv.Itoa(int(e.ErrorCode)) + ")"
}

func (e FondyError) CodeIs(code fondy_status.StatusCode) bool {
	return e.ErrorCode == code
}

func (e FondyError) IsFatalError() bool {
	return e.IsFatal
}

func (e FondyError) Code() fondy_status.StatusCode {
	return e.ErrorCode
}
