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
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/stremovskyy/gofondy/recorder"

	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/models"
	"github.com/stremovskyy/gofondy/models/models_v2"
	"github.com/stremovskyy/gofondy/utils"
)

type v2Client struct {
	client   *http.Client
	logger   *log.Logger
	recorder recorder.Client
}

func (m *v2Client) do(url consts.FondyURL, order *models_v2.Order, credit bool, merchantAccount *models.MerchantAccount, addOrderDescription bool) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	if addOrderDescription {
		order.OrderDesc = utils.StringRef(merchantAccount.MerchantString)
	}

	wholeAmount, err := strconv.ParseFloat(*order.Amount, 64)
	if err != nil {
		return nil, errors.New("split accounts problem: amount parse error")
	}

	splitAmountSum := 0.0

	for _, splitAccount := range merchantAccount.SplitAccounts {
		splitAmount := wholeAmount * splitAccount.SplitPercentage / 100
		merchantReceiver := models_v2.NewMerchantReceiver(models_v2.NewMerchantRequisites(int64(splitAmount), &splitAccount.MerchantID, &splitAccount.MerchantAddedDescription))
		order.Receiver = append(order.Receiver, *merchantReceiver)
		splitAmountSum += splitAmount
	}

	if splitAmountSum != wholeAmount {
		return nil, fmt.Errorf("order %s split accounts problem: split amount sum %f != whole amount %f", *order.OrderID, splitAmountSum, wholeAmount)
	}

	fondyRequest := models_v2.NewRequest(order)

	if credit {
		fondyRequest.Sign(merchantAccount.MerchantCreditKey)
	} else {
		fondyRequest.Sign(merchantAccount.MerchantKey)
	}

	jsonValue, err := json.Marshal(fondyRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	ctx := context.WithValue(context.Background(), "request_id", requestID)
	tags := tagsOrderRetriever(order)

	if m.recorder != nil {
		err = m.recorder.RecordRequest(ctx, nil, requestID, jsonValue, tags)
		if err != nil {
			m.logger.Printf("[ERROR] cannot record request: %v", err)
		}
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

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		if m.recorder != nil {
			err = m.recorder.RecordError(ctx, nil, requestID, err, tags)
			if err != nil {
				m.logger.Printf("[ERROR] cannot record request error: %v", err)
			}
		}

		return nil, fmt.Errorf("cannot copy response buffer: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}(resp.Body)

	if m.recorder != nil {
		err = m.recorder.RecordResponse(ctx, nil, requestID, raw, tags)
		if err != nil {
			m.logger.Printf("[ERROR] cannot record response: %v", err)
		}
	}

	errorResponse, _ := models_v2.UnmarshalErrorResponse(raw)
	if errorResponse.Response.ErrorCode != 0 {
		return nil, fmt.Errorf("fondy error response (%d): %s", errorResponse.Response.ErrorCode, errorResponse.Response.ErrorMessage)
	}

	return &raw, nil
}

func tagsOrderRetriever(order *models_v2.Order) map[string]string {
	tags := make(map[string]string)

	if order.OrderID != nil {
		tags["order_id"] = *order.OrderID
	}

	return tags
}
