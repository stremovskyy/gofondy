/*
 * MIT License
 *
 * Copyright (c) 2024 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
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

type idClient struct {
	client   *http.Client
	options  *ClientOptions
	logger   *log.Logger
	recorder recorder.Client
}

func (c *idClient) clientStatus(fondyURL consts.FondyURL, request *models.FondyClientStatusRequest) (*[]byte, error) {
	// Generate a unique request ID
	requestID := uuid.New().String()
	methodPost := "POST"
	ctx := context.WithValue(context.Background(), "request_id", requestID)
	tags := tagsRetriever(request)

	metricsMap := make(map[string]string)
	metricsMap["url"] = fondyURL.String()
	tim := time.Now()
	metricsMap["start_timestamp"] = tim.Format("2006-01-02 15:04:05")
	defer func() {
		metricsMap["end_timestamp"] = time.Now().Format("2006-01-02 15:04:05")
		metricsMap["duration"] = fmt.Sprintf("%s", time.Since(tim).String())

		if c.recorder != nil {
			err := c.recorder.RecordMetrics(ctx, nil, requestID, metricsMap, tags)
			if err != nil {
				c.logger.Printf("[ERROR] cannot record metrics: %v", err)
			}
		}
	}()

	// Debug logging
	if c.options.IsDebug {
		c.logger.Printf("[GO FONDY] Request ID: %v\n", requestID)
		c.logger.Printf("[GO FONDY] URL: %v\n", fondyURL.String())
		c.logger.Printf("[GO FONDY] Client Status Request: %+v\n", request)
	}

	// Serialize the request object to JSON
	jsonValue, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	if c.options.IsDebug {
		c.logger.Printf("[GO FONDY] Request: %v\n", string(jsonValue))
	}

	if c.recorder != nil {
		err = c.recorder.RecordRequest(ctx, nil, requestID, jsonValue, tags)
		if err != nil {
			c.logger.Printf("[ERROR] cannot record request: %v", err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest(methodPost, fondyURL.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	// Set the necessary headers here
	req.Header.Set("User-Agent", "GOFONDY/"+consts.Version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("X-API-Version", "1.0")

	// Send the request using the client's http.Client
	resp, err := c.client.Do(req)
	if err != nil {
		if c.recorder != nil {
			err = c.recorder.RecordError(ctx, nil, requestID, err, tags)
			if err != nil {
				c.logger.Printf("[ERROR] cannot record request error: %v", err)
			}
		}

		return nil, fmt.Errorf("cannot send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	if c.recorder != nil {
		err = c.recorder.RecordResponse(ctx, nil, requestID, responseBody, tags)
		if err != nil {
			c.logger.Printf("[ERROR] cannot record response: %v", err)
		}
	}

	// Debug logging for the response
	if c.options.IsDebug {
		c.logger.Printf("[GO FONDY] Response: %v\n", string(responseBody))
	}

	return &responseBody, nil
}

func tagsRetriever(request *models.FondyClientStatusRequest) map[string]string {
	tags := make(map[string]string)

	if request.MerchantID != nil {
		tags["merchant:id"] = *request.MerchantID
	}

	if request.IPN != nil {
		tags["id:ipn"] = *request.IPN
	}

	if request.IDCard != nil {
		tags["id:card"] = *request.IDCard
	}

	if request.InternalPassport != nil {
		tags["id:internal_passport"] = *request.InternalPassport
	}

	return tags
}
