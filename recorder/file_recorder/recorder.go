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

package file_recorder

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stremovskyy/gofondy/recorder"
)

type fileRecorder struct {
	requestLogger  *log.Logger
	responseLogger *log.Logger
	errorLogger    *log.Logger
}

func (r *fileRecorder) RecordMetrics(ctx context.Context, orderID *string, requestID string, metrics map[string]string, tags map[string]string) error {
	// File-based logs do not support recording metrics.
	return fmt.Errorf("recordMetrics not supported in file-based recorder")
}

// NewFileRecorder creates a new instance of fileRecorder.
func NewFileRecorder(requestLogPath, responseLogPath, errorLogPath string) recorder.Client {
	requestLogFile, err := os.OpenFile(requestLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open request log file: " + err.Error())
	}

	responseLogFile, err := os.OpenFile(responseLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open response log file: " + err.Error())
	}

	errorLogFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open error log file: " + err.Error())
	}

	return &fileRecorder{
		requestLogger:  log.New(requestLogFile, "REQUEST: ", log.LstdFlags),
		responseLogger: log.New(responseLogFile, "RESPONSE: ", log.LstdFlags),
		errorLogger:    log.New(errorLogFile, "ERROR: ", log.LstdFlags),
	}
}

func (r *fileRecorder) RecordRequest(ctx context.Context, id *string, requestID string, request []byte, tags map[string]string) error {
	logEntry := fmt.Sprintf("%s - RequestID: %s, Tags: %v, Request: %s", time.Now().Format(time.RFC3339), requestID, tags, string(request))
	r.requestLogger.Println(logEntry)
	return nil
}

func (r *fileRecorder) RecordResponse(ctx context.Context, id *string, requestID string, response []byte, tags map[string]string) error {
	logEntry := fmt.Sprintf("%s - RequestID: %s, Tags: %v, Response: %s", time.Now().Format(time.RFC3339), requestID, tags, string(response))
	r.responseLogger.Println(logEntry)
	return nil
}

func (r *fileRecorder) RecordError(ctx context.Context, id *string, requestID string, err error, tags map[string]string) error {
	logEntry := fmt.Sprintf("%s - RequestID: %s, Tags: %v, Error: %s", time.Now().Format(time.RFC3339), requestID, tags, err.Error())
	r.errorLogger.Println(logEntry)
	return nil
}

func (r *fileRecorder) GetRequest(ctx context.Context, requestID string) ([]byte, error) {
	// File-based logs do not support retrieving specific entries.
	return nil, fmt.Errorf("getRequest not supported in file-based recorder")
}

func (r *fileRecorder) GetResponse(ctx context.Context, requestID string) ([]byte, error) {
	// File-based logs do not support retrieving specific entries.
	return nil, fmt.Errorf("getResponse not supported in file-based recorder")
}

func (r *fileRecorder) FindByTag(ctx context.Context, tag string) ([]string, error) {
	// File-based logs do not support querying by tags.
	return nil, fmt.Errorf("findByTag not supported in file-based recorder")
}
