/*
 * Project: banker
 * File: options.go (2/20/24, 11:23 AM)
 *
 * Copyright (C) Megakit Systems 2017-2024, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package redis_recorder

import "time"

type Options struct {
	Debug          bool
	Addr           string
	Password       string
	DB             int
	DefaultTTL     time.Duration
	CompressionLvl int
	Prefix         string
}

func NewDefaultOptions(addr string, password string, DB int) *Options {
	return &Options{
		Addr:           addr,
		Password:       password,
		DB:             DB,
		DefaultTTL:     24 * time.Hour * 7, // one week
		CompressionLvl: 3,
		Prefix:         "fondy:http",
	}
}
