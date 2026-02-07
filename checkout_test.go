package sepay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

func TestCheckoutService_InitCheckoutURL(t *testing.T) {
	t.Run("sandbox", func(t *testing.T) {
		c, _ := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
		})
		expected := "https://pay-sandbox.sepay.vn/v1/checkout/init"
		if got := c.Checkout.InitCheckoutURL(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("production", func(t *testing.T) {
		c, _ := NewClient(Config{
			Env:        Production,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
		})
		expected := "https://pay.sepay.vn/v1/checkout/init"
		if got := c.Checkout.InitCheckoutURL(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func TestCheckoutService_InitOneTimePaymentFields(t *testing.T) {
	t.Run("minimal fields with defaults", func(t *testing.T) {
		c, _ := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "test_merchant",
			SecretKey:  "test_secret",
		})

		signed := c.Checkout.InitOneTimePaymentFields(OnetimePaymentFields{
			OrderInvoiceNumber: "INV-001",
			OrderAmount:        50000,
			Currency:           "VND",
			OrderDescription:   "Test order",
		})

		if signed.Merchant != "test_merchant" {
			t.Errorf("expected merchant %q, got %q", "test_merchant", signed.Merchant)
		}
		if signed.Operation != OperationPurchase {
			t.Errorf("expected operation %q, got %q", OperationPurchase, signed.Operation)
		}
		if signed.Signature == "" {
			t.Error("expected non-empty signature")
		}

		// Verify signature manually.
		data := strings.Join([]string{
			"merchant=test_merchant",
			"operation=PURCHASE",
			"order_amount=50000",
			"currency=VND",
			"order_invoice_number=INV-001",
			"order_description=Test order",
		}, ",")
		mac := hmac.New(sha256.New, []byte("test_secret"))
		mac.Write([]byte(data))
		expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		if signed.Signature != expectedSig {
			t.Errorf("signature mismatch:\n  expected: %s\n  got:      %s", expectedSig, signed.Signature)
		}
	})

	t.Run("all optional fields", func(t *testing.T) {
		c, _ := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "test_merchant",
			SecretKey:  "test_secret",
		})

		signed := c.Checkout.InitOneTimePaymentFields(OnetimePaymentFields{
			Operation:          OperationPurchase,
			PaymentMethod:      BankTransfer,
			OrderInvoiceNumber: "INV-002",
			OrderAmount:        100000.5,
			Currency:           "VND",
			OrderDescription:   "Full test",
			OrderTaxAmount:     Float64(1000),
			CustomerID:         String("CUST-001"),
			SuccessURL:         String("https://example.com/success"),
			ErrorURL:           String("https://example.com/error"),
			CancelURL:          String("https://example.com/cancel"),
			CustomData:         String("extra-data"),
		})

		if signed.PaymentMethod != BankTransfer {
			t.Errorf("expected payment_method %q, got %q", BankTransfer, signed.PaymentMethod)
		}
		if signed.CustomerID == nil || *signed.CustomerID != "CUST-001" {
			t.Errorf("expected customer_id %q", "CUST-001")
		}
		if signed.CustomData == nil || *signed.CustomData != "extra-data" {
			t.Errorf("expected custom_data %q", "extra-data")
		}

		// Verify signature includes payment_method, customer_id, and URLs.
		data := strings.Join([]string{
			"merchant=test_merchant",
			"operation=PURCHASE",
			"payment_method=BANK_TRANSFER",
			"order_amount=100000.5",
			"currency=VND",
			"order_invoice_number=INV-002",
			"order_description=Full test",
			"customer_id=CUST-001",
			"success_url=https://example.com/success",
			"error_url=https://example.com/error",
			"cancel_url=https://example.com/cancel",
		}, ",")
		mac := hmac.New(sha256.New, []byte("test_secret"))
		mac.Write([]byte(data))
		expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		if signed.Signature != expectedSig {
			t.Errorf("signature mismatch:\n  expected: %s\n  got:      %s", expectedSig, signed.Signature)
		}
	})
}

func TestSignedCheckoutFields_FormValues(t *testing.T) {
	c, _ := NewClient(Config{
		Env:        Sandbox,
		MerchantID: "test_merchant",
		SecretKey:  "test_secret",
	})

	signed := c.Checkout.InitOneTimePaymentFields(OnetimePaymentFields{
		OrderInvoiceNumber: "INV-001",
		OrderAmount:        50000,
		Currency:           "VND",
		OrderDescription:   "Test order",
		CustomerID:         String("CUST-001"),
	})

	form := signed.FormValues()

	expected := map[string]string{
		"merchant":             "test_merchant",
		"operation":            "PURCHASE",
		"order_invoice_number": "INV-001",
		"order_amount":         "50000",
		"currency":             "VND",
		"order_description":    "Test order",
		"customer_id":          "CUST-001",
		"signature":            signed.Signature,
	}

	for k, v := range expected {
		if form[k] != v {
			t.Errorf("form[%q] = %q, want %q", k, form[k], v)
		}
	}

	if _, ok := form["payment_method"]; ok {
		t.Error("expected payment_method to be absent")
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{10000, "10000"},
		{100.5, "100.5"},
		{0, "0"},
		{1.23456, "1.23456"},
		{100000.00, "100000"},
	}
	for _, tc := range tests {
		got := formatFloat(tc.input)
		if got != tc.expected {
			t.Errorf("formatFloat(%v) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
