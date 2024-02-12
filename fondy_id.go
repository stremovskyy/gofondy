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

package gofondy

import (
	"encoding/json"
	"fmt"

	"github.com/stremovskyy/gofondy/manager"
	"github.com/stremovskyy/gofondy/models"
)

var id ID

type fondyID struct {
	manager manager.FondyManager
	options *models.Options
}

func (f *fondyID) Status(statusRequest *models.IDStatusRequest) (*models.FondyClientStatusResponse, error) {

	fondyStatusRequest := &models.FondyClientStatusRequest{
		MerchantID: statusRequest.GetMerchantID(),
	}

	switch statusRequest.IDType {
	case models.IDTypeTIN:
		fondyStatusRequest.IPN = statusRequest.IDref()
	case models.IDTypePassport:
		fondyStatusRequest.InternalPassport = statusRequest.IDref()
	case models.IDTypeIDCard:
		fondyStatusRequest.IDCard = statusRequest.IDref()
	}

	err := fondyStatusRequest.Sign(statusRequest.Merchant.MerchantKey, f.options.IsDebug)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	rawStatusResponse, err := f.manager.IDStatus(fondyStatusRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var response models.FondyClientStatusResponse

	err = json.Unmarshal(*rawStatusResponse, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.IsError() {
		return nil, fmt.Errorf("response error: %w", response.GetError())
	}

	return &response, nil
}

func (f *fondyID) Limits(limitsRequest *models.IDStatusRequest) (*models.FondyBalance, error) {
	wholeResponse, err := f.Status(limitsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return wholeResponse.Balance, nil
}

func (g *gateway) ID() ID {
	if id == nil {
		id = &fondyID{
			manager: g.manager,
			options: g.options,
		}
	}

	return id
}
