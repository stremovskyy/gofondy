/*
 * Project: banker
 * File: fondy_client_status_response.go (2/19/24, 5:13 PM)
 *
 * Copyright (C) Megakit Systems 2017-2024, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package models

import (
	"encoding/json"
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

func (r *FondyClientStatusResponse) Bytes() []byte {
	jsonString, _ := json.Marshal(r)

	return jsonString
}
