# Gofondy - Fondy Payment Gate Client 

[![GoDoc](https://godoc.org/github.com/karmadon/gofondy?status.svg)](https://godoc.org/github.com/karmadon/gofondy)

**gofondy** is GO client for the Fondy Payment Gate API.  

## Jump to

* [Installation](#Installation)
* [Usage](#Usage)
* [API examples](#API-examples)
* [License](#License)
* [Contributing](#Contributing)
* [Authors](#Authors)
* [Acknowledgments](#Acknowledgments)
* [TODO](#TODO)

## Installation

```bash
go get github.com/karmadon/gofondy
```

## Usage

### Create client
create a new Fondy client with default options

```go
fondyGateway := gofondy.New(models.DefaultOptions())
```

### Using a merchant account
Merchant account is a structure that contains all the necessary information for the client to work with the Fondy API.

```go
merchAccount := &models.MerchantAccount{
    MerchantID:       examples.TechMerchantId,
    MerchantKey:      examples.TechMerchantKey,
    MerchantString:   "Test Merchant",
    MerchantDesignID: examples.DesignId,
    IsTechnical:      true,
}
```

## API examples
### Card verification

```go
verificationLink, err := fondyGateway.VerificationLink(merchAccount, uuid.New(), nil, "test", consts.CurrencyCodeUAH)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("\nVerification link: %s\n", verificationLink.String())
```

### Payment Hold

```go
invoiceId := uuid.New()

holdAmount := float64(3)

paymentByToken, err := fondyGateway.Hold(merchAccount, &invoiceId, &holdAmount, examples.CardToken)
if err != nil {
    log.Fatal(err)
}

if *paymentByToken.ResponseStatus == consts.FondyResponseStatusSuccess {
    fmt.Printf("Order (%s) status: %s\n", paymentByToken.OrderID, *paymentByToken.OrderStatus)
} else {
    fmt.Printf("Error: %s\n", *paymentByToken.ErrorMessage)
}
```

### Payment Capture

```go
invoiceId := uuid.MustParse("767f44ef-2997-4623-961f-9ee081ef730f")

captureAmount := float64(3)

capturePayment, err := fondyGateway.Capture(merchAccount, &invoiceId, &captureAmount)
if err != nil {
    log.Fatal(err)
}

if *capturePayment.ResponseStatus == consts.FondyResponseStatusSuccess {
    fmt.Printf("Order (%s) Capture Status: %s\n", capturePayment.OrderID.String(), *capturePayment.CaptureStatus)
} else {
    fmt.Printf("Error: %s\n", *capturePayment.ErrorMessage)
}
```

### Payment Refund

```go
	refundPayment, err := fondyGateway.Refund(merchAccount, &invoiceId, &captureAmount)
	if err != nil {
		log.Fatal(err)
	}

	if *refundPayment.ResponseStatus == consts.FondyResponseStatusSuccess {
		fmt.Printf("Order (%s) Reversed: %s\n", refundPayment.OrderID.String(), *refundPayment.ReversalAmount)
	} else {
		fmt.Printf("Error: %s\n", *refundPayment.ErrorMessage)
	}
```

### Payment Status

```go
	status, err := fondyGateway.Status(merchAccount, &invoiceId)
	if err != nil {
		log.Fatal(err)
	}

	if *status.ResponseStatus == consts.FondyResponseStatusSuccess {
		fmt.Printf("Order status: %s\n", *status.OrderStatus)
	} else {
		fmt.Printf("Error: %s\n", *status.ErrorMessage)
	}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments
- [Fondy API Documentation](https://docs.fondy.eu/en/)
- [GoDoc](https://godoc.org/github.com/karmadon/gofondy)

## Author

* **Anton Stremovskyy** - *Initial work* - [Karmadon](https://github.com/karmadon)

### Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
Please make sure to update tests as appropriate.

## TODO
- [ ] Add tests
- [ ] Add more examples
- [ ] Add more API methods
