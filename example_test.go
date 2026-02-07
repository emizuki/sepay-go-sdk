package sepay_test

import (
	"context"
	"fmt"
	"log"

	"github.com/emizuki/sepay-go-sdk"
)

func ExampleNewClient() {
	client, err := sepay.NewClient(sepay.Config{
		Env:        sepay.Sandbox,
		MerchantID: "your_merchant_id",
		SecretKey:  "your_secret_key",
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = client
}

func ExampleOrderService_All() {
	client, err := sepay.NewClient(sepay.Config{
		Env:        sepay.Sandbox,
		MerchantID: "your_merchant_id",
		SecretKey:  "your_secret_key",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Order.All(context.Background(), &sepay.OrderQueryParams{
		PerPage:     sepay.Int(10),
		OrderStatus: sepay.String("COMPLETED"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
}

func ExampleCheckoutService_InitOneTimePaymentFields() {
	client, err := sepay.NewClient(sepay.Config{
		Env:        sepay.Sandbox,
		MerchantID: "your_merchant_id",
		SecretKey:  "your_secret_key",
	})
	if err != nil {
		log.Fatal(err)
	}

	checkoutURL := client.Checkout.InitCheckoutURL()
	signed := client.Checkout.InitOneTimePaymentFields(sepay.OnetimePaymentFields{
		OrderInvoiceNumber: "INV-001",
		OrderAmount:        50000,
		Currency:           "VND",
		OrderDescription:   "Payment for Order INV-001",
		PaymentMethod:      sepay.BankTransfer,
		SuccessURL:         sepay.String("https://example.com/success"),
		ErrorURL:           sepay.String("https://example.com/error"),
		CancelURL:          sepay.String("https://example.com/cancel"),
	})

	fmt.Println("Checkout URL:", checkoutURL)
	fmt.Println("Signature:", signed.Signature)

	// Use signed.FormValues() to build an HTML form.
	for key, value := range signed.FormValues() {
		fmt.Printf("<input type=\"hidden\" name=\"%s\" value=\"%s\">\n", key, value)
	}
}
