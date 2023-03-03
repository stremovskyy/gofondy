/*
 * MIT License
 *
 * Copyright (c) 2022 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
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

package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
)

type v1Client struct {
	client *http.Client
}

func (m *v1Client) do(url consts.FondyURL, request *models.RequestObject, credit bool, merchantAccount *models.MerchantAccount) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

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
		"User-Agent":    {"GOFONDY/" + consts.Version},
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
		"X-Request-ID":  {requestID},
		"X-API-Version": {"1.0"},
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
