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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/stremovskyy/gofondy/recorder"

	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
)

type v1Client struct {
	client   *http.Client
	options  *ClientOptions
	logger   *log.Logger
	recorder recorder.Client
}

func (m *v1Client) do(url consts.FondyURL, request *models.FondyRequestObject, credit bool, merchantAccount *models.MerchantAccount, reservationData *models.ReservationData) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	tags := tagsRequestRetriever(request)
	ctx := context.WithValue(context.Background(), "request_id", requestID)

	if reservationData != nil {
		request.ReservationData = reservationData.Base64Encoded()
	}

	if m.options.IsDebug {
		m.logger.Printf("[GO FONDY] Request ID: %v\n", requestID)
		m.logger.Printf("[GO FONDY] URL: %v\n", url.String())
		m.logger.Printf("[GO FONDY] Reservation data: %v\n", reservationData)
	}

	metricsMap := make(map[string]string)
	metricsMap["url"] = url.String()
	tim := time.Now()
	metricsMap["start_timestamp"] = tim.Format("2006-01-02 15:04:05")
	defer func() {
		metricsMap["end_timestamp"] = tim.Format("2006-01-02 15:04:05")
		metricsMap["duration"] = fmt.Sprintf("%s", time.Since(tim).String())

		if m.recorder != nil {
			err := m.recorder.RecordMetrics(ctx, nil, requestID, metricsMap, tags)
			if err != nil {
				m.logger.Printf("[ERROR] cannot record metrics: %v", err)
			}
		}
	}()

	var key string
	if credit {
		key = merchantAccount.MerchantCreditKey
	} else {
		key = merchantAccount.MerchantKey
	}

	err := request.Sign(key, m.options.IsDebug)
	if err != nil {
		return nil, fmt.Errorf("cannot sign request: %v", err)
	}

	jsonValue, err := json.Marshal(models.NewFondyRequest(request))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	if m.options.IsDebug {
		m.logger.Printf("[GO FONDY] Request: %v\n", string(jsonValue))
	}

	if m.recorder != nil {
		err = m.recorder.RecordRequest(ctx, request.OrderID, requestID, jsonValue, tags)
		if err != nil {
			m.logger.Printf("[ERROR] cannot record request: %v", err)
		}
	}

	req, err := http.NewRequest(methodPost, url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("User-Agent", "GOFONDY/"+consts.Version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("X-API-Version", "1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		if m.recorder != nil {
			err = m.recorder.RecordError(ctx, request.OrderID, requestID, err, tags)
			if err != nil {
				m.logger.Printf("[ERROR] cannot record request error: %v", err)
			}
		}

		return nil, fmt.Errorf("cannot send request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	if m.recorder != nil {
		err = m.recorder.RecordResponse(ctx, request.OrderID, requestID, raw, tags)
		if err != nil {
			m.logger.Printf("[ERROR] cannot record response: %v", err)
		}
	}

	if m.options.IsDebug {
		log.Printf("[GO FONDY] Response: %v\n", string(raw))
	}

	return &raw, nil
}

func tagsRequestRetriever(request *models.FondyRequestObject) map[string]string {
	tags := make(map[string]string)

	if request.OrderID != nil {
		tags["order_id"] = *request.OrderID
	}

	return tags
}
