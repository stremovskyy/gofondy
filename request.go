package gofondy

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Request struct {
	Request RequestObject `json:"request"`
}

func NewFondyRequest(request RequestObject) *Request {
	return &Request{Request: request}
}

// RequestObject Accept purchase (hosted payment page)
type RequestObject struct {
	OrderID    *string `json:"order_id"`
	MerchantID *string `json:"merchant_id"`
	Signature  *string `json:"signature"`

	OrderDesc         *string `json:"order_desc,omitempty"`
	Amount            *string `json:"amount,omitempty"`
	Currency          *string `json:"currency,omitempty"`
	Preauth           *string `json:"preauth,omitempty"`
	DesignID          *string `json:"design_id,omitempty"`
	Rectoken          *string `json:"rectoken,omitempty"`
	ProductID         *string `json:"product_id,omitempty"`
	Lang              *string `json:"lang,omitempty"`
	SenderEmail       *string `json:"sender_email,omitempty"`
	ServerCallbackURL *string `json:"server_callback_url,omitempty"`
	Lifetime          *string `json:"lifetime,omitempty"`
	Verification      *string `json:"verification,omitempty"`
	RequiredRectoken  *string `json:"required_rectoken,omitempty"`
	MerchantData      *string `json:"merchant_data,omitempty"`
	ReceiverRectoken  *string `json:"receiver_rectoken"`
}

// CreateSignature - creates signature for fondy payment gate
func (r *RequestObject) CreateSignature(merchantKey string) error {
	s := ""

	s += merchantKey + "|"

	values := reflect.ValueOf(*r)
	types := values.Type()
	preFiltered := map[string]string{}

	for i := 0; i < values.NumField(); i++ {
		t := values.Field(i).Interface()
		if t != nil {
			s, ok := t.(*string)
			if ok && s != nil {
				preFiltered[types.Field(i).Name] = *s
			}
		}
	}

	keys := make([]string, 0, len(preFiltered))
	for k := range preFiltered {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	final := make([]string, 0, len(preFiltered))
	for _, k := range keys {
		final = append(final, preFiltered[k])
	}

	s += strings.Join(final, "|")

	h := sha1.New()
	h.Write([]byte(s))
	r.Signature = StringRef(fmt.Sprintf("%x", h.Sum(nil)))

	return nil
}
