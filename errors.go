package gofondy

import (
	"strconv"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`

	RequestObject *RequestObject `json:"-"`
	RawResponse   *[]byte        `json:"-"`
}

func NewAPIError(code int, message string, err error, requestObject *RequestObject, rawResponse *[]byte) *APIError {
	return &APIError{Code: code, Message: message, Err: err, RequestObject: requestObject, RawResponse: rawResponse}
}

func (e APIError) Error() string {
	return "HTTP error: " + e.Message + " (" + strconv.Itoa(e.Code) + ")" + " " + e.Err.Error()
}
