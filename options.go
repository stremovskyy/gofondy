package gofondy

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

	CallbackUrl string
	DesignId    string
	MerchantId  string
	MerchantKey string
}
