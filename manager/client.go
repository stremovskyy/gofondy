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
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy/consts"
	"github.com/karmadon/gofondy/models"
	"github.com/karmadon/gofondy/models/models_v2"
	"github.com/karmadon/gofondy/utils"
)

func (m *manager) payment(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, true)
}

func (m *manager) splitPayment(url consts.FondyURL, request *models_v2.SplitRequest, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.doWithSplit(url, request, false, merchantAccount, true)
}

func (m *manager) verify(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, true)
}

func (m *manager) withdraw(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	return m.do(url, request, true, merchantAccount, true)
}

func (m *manager) final(url consts.FondyURL, request *models.RequestObject, merchantAccount *models.MerchantAccount, technical bool) (*[]byte, error) {
	return m.do(url, request, false, merchantAccount, false)
}

func (m *manager) do(url consts.FondyURL, request *models.RequestObject, credit bool, merchantAccount *models.MerchantAccount, addOrderDescription bool) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	request.MerchantID = &merchantAccount.MerchantID

	if addOrderDescription {
		request.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	}

	if credit {
		err := request.Sign(merchantAccount.MerchantCreditKey)
		if err != nil {
			return nil, fmt.Errorf("cannot sign request with credit key: %v", err)
		}
	} else {
		err := request.Sign(merchantAccount.MerchantKey)
		if err != nil {
			return nil, fmt.Errorf("cannot sign request with merchant key: %v", err)
		}
	}

	jsonValue, err := json.Marshal(models.NewFondyRequest(request))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	req, err := http.NewRequest(methodPost, url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header = http.Header{
		"User-Agent":   {"GOFONDY/" + consts.Version},
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"X-Request-ID": {requestID},
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot copy response buffer: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}(resp.Body)

	return &raw, nil
}

func (m *manager) doWithSplit(url consts.FondyURL, request *models_v2.SplitRequest, credit bool, merchantAccount *models.MerchantAccount, addOrderDescription bool) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	merchantId, err := strconv.ParseInt(merchantAccount.MerchantID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse merchant id: %w", err)
	}

	request.Order.MerchantID = &merchantId

	if addOrderDescription {
		request.Order.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	}

	fondyRequest := models_v2.NewRequest(request)

	if credit {
		fondyRequest.Sign(merchantAccount.MerchantCreditKey)
	} else {
		fondyRequest.Sign(merchantAccount.MerchantKey)
	}

	jsonValue, err := json.Marshal(fondyRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	req, err := http.NewRequest(methodPost, url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header = http.Header{
		"User-Agent":      {"GOFONDY/" + consts.Version},
		"Accept":          {"application/json"},
		"Content-Type":    {"application/json"},
		"Accept-Encoding": {"gzip"},
		"X-Request-ID":    {requestID},
		"X-API-Version":   {"2.0"},
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				log.Printf("cannot close gzip reader: %v", err)
			}
		}(reader)
	default:
		reader = resp.Body
	}

	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return nil, fmt.Errorf("cannot copy response buffer: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}(resp.Body)

	return &raw, nil
}
