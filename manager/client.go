/*
 * Project: banker
 * File: client.go (8/22/22, 2:57 PM)
 *
 * Copyright (C) Megakit Systems 2017-2022, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (karmadon) Stremovskyy <stremovskyy@gmail.com>
 */

package manager

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/karmadon/gofondy/consts"
	"github.com/karmadon/gofondy/models"
	"github.com/karmadon/gofondy/models/models_v2"
)

type Client interface {
	payment(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
	split(url consts.FondyURL, order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error)
	withdraw(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error)
}

type client struct {
	v1 *v1Client
	v2 *v2Client
}

type ClientOptions struct {
	Timeout         time.Duration
	KeepAlive       time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
}

func NewClient(options *ClientOptions) Client {
	dialer := &net.Dialer{
		Timeout:   options.Timeout,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}}

	cl := &http.Client{Transport: tr}

	return &client{
		v1: &v1Client{
			client: cl,
		},
		v2: &v2Client{
			client: cl,
		},
	}
}

func (m *client) payment(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.v1.do(url, request, false, merchantAccount)
}

func (m *client) withdraw(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.v1.do(url, request, true, merchantAccount)
}

func (m *client) split(url consts.FondyURL, order *models_v2.Order, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.v2.do(url, order, false, merchantAccount, true)
}
