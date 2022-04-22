/*
 * Project: go-driveapp-ms
 * File: fondy_verification.go (8/8/19, 6:52 PM)
 *
 * Copyright (C) Megakit Systems 2017-2019, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (karmadon) Stremovskyy <stremovskyy@gmail.com>
 */

package gofondy

type VerificationResponceResponse struct {
	ResponseStatus string `json:"response_status"`
	ErrorMessage   string `json:"error_message"`
	ErrorCode      int64  `json:"error_code"`
	RequestID      string `json:"request_id"`
}
