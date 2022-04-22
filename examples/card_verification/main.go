package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/karmadon/gofondy"
)

func main() {
	options := &gofondy.Options{
		Timeout:                 30 * time.Second,
		KeepAlive:               30 * time.Second,
		MaxIdleConns:            10,
		IdleConnTimeout:         20 * time.Second,
		VerificationAmount:      1,
		VerificationDescription: "Verification Test",
		VerificationLifeTime:    600 * time.Second,
		CallbackUrl:             FondyVerificationCallbackURL,
		DesignId:                DesignId,
		MerchantId:              MerchantId,
		MerchantKey:             MerchantKey,
	}

	fondyGateway := gofondy.New(options)

	verificationLink, err := fondyGateway.VerificationLink(uuid.New(), nil, "test", gofondy.CurrencyCodeUAH)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Verification link: %s", *verificationLink)
}
