package gofondy

import (
	"github.com/google/uuid"
)

type FondyGateway interface {
	VerificationLink(invoiceId uuid.UUID, email *string, note string, code CurrencyCode) (*string, error)
}
