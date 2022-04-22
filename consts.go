package gofondy

const Version = "0.0.1"

const (
	FondyTimeFormat = "02.01.2006 15:04:05"
)

type FondyURL string

const (
	FondyURLGetVerification FondyURL = "https://api.fondy.eu/api/checkout/url/"
)

func (t FondyURL) String() string {
	return string(t)
}

type FondyTransactionType string

const (
	FondyTransactionTypePurchase     FondyTransactionType = "purchase"
	FondyTransactionTypeReverse      FondyTransactionType = "reverse"
	FondyTransactionTypeVerification FondyTransactionType = "verification"
	FondyTransactionTypeP2PCredit    FondyTransactionType = "p2p credit"
)

func (t FondyTransactionType) String() string {
	return string(t)
}

type FondyResponseStatus string

const (
	FondyResponseStatusSuccess FondyResponseStatus = "success"
	FondyResponseStatusFailure FondyResponseStatus = "failure"
)

func (s FondyResponseStatus) String() string {
	return string(s)
}

type CurrencyCode string

const (
	CurrencyCodeUAH CurrencyCode = "UAH"
)

func (c CurrencyCode) String() string {
	return string(c)
}

type FondyCardType string

const (
	FondyCardTypeVISA       FondyCardType = "VISA"
	FondyCardTypeMASTERCARD FondyCardType = "MASTERCARD"
)

func (t FondyCardType) String() string {
	return string(t)
}

type FondyCaptureStatus string

const (
	FondyCaptureStatusHold     FondyCaptureStatus = "hold"
	FondyCaptureStatusCaptured FondyCaptureStatus = "captured"
)

func (s *FondyCaptureStatus) String() string {
	if s != nil {
		return string(*s)
	}
	return ""
}

// Reversal processing status
type FondyReverseStatus string

const (
	// reversal has been created, but not processed yet
	FondyReverseStatusCreated FondyReverseStatus = "created"
	// reversal is declined by FONDY payment gateway or by bank or by external payment system
	FondyReverseStatusDeclined FondyReverseStatus = "declined"
	// reversal completed successfully
	FondyReverseStatusApproved FondyReverseStatus = "approved"
)

type Status string

const (
	StatusReversed   Status = "reversed"
	StatusApproved   Status = "approved"
	StatusProcessing Status = "processing"
	StatusDeclined   Status = "declined"
	StatusExpired    Status = "expired"
	StatusCreated    Status = "created"
	StatusCanceled   Status = "canceled"
	StatusCaptured   Status = "captured"
)
