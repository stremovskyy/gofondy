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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
)

type idClient struct {
	client  *http.Client
	options *ClientOptions
}

func (c *idClient) clientStatus(fondyURL consts.FondyURL, request *models.FondyClientStatusRequest) (*[]byte, error) {
	// Generate a unique request ID
	requestID := uuid.New().String()
	methodPost := "POST"

	// Debug logging
	if c.options.IsDebug {
		log.Printf("[GO FONDY] Request ID: %v\n", requestID)
		log.Printf("[GO FONDY] URL: %v\n", fondyURL.String())
		log.Printf("[GO FONDY] Client Status Request: %+v\n", request)
	}

	// Serialize the request object to JSON
	jsonValue, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	// More debug logging
	if c.options.IsDebug {
		log.Printf("[GO FONDY] JSON Request: %v\n", string(jsonValue))
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

	// Debug logging for the response
	if c.options.IsDebug {
		log.Printf("[GO FONDY] Response: %v\n", string(responseBody))
	}

	return &responseBody, nil
}
