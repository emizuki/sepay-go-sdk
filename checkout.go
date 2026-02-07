package sepay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
)

// PaymentMethod represents a supported payment method.
type PaymentMethod string

const (
	Card              PaymentMethod = "CARD"
	BankTransfer      PaymentMethod = "BANK_TRANSFER"
	NapasBankTransfer PaymentMethod = "NAPAS_BANK_TRANSFER"
)

// Operation represents a checkout operation type.
type Operation string

const (
	OperationPurchase Operation = "PURCHASE"
	OperationVerify   Operation = "VERIFY"
)

// OnetimePaymentFields holds the fields for a one-time payment checkout.
type OnetimePaymentFields struct {
	Operation          Operation
	PaymentMethod      PaymentMethod
	OrderInvoiceNumber string
	OrderAmount        float64
	Currency           string
	OrderDescription   string
	OrderTaxAmount     *float64
	CustomerID         *string
	SuccessURL         *string
	ErrorURL           *string
	CancelURL          *string
	CustomData         *string
}

// SignedCheckoutFields holds all checkout fields including the computed signature.
type SignedCheckoutFields struct {
	Merchant           string        `json:"merchant"`
	Operation          Operation     `json:"operation"`
	PaymentMethod      PaymentMethod `json:"payment_method,omitempty"`
	OrderInvoiceNumber string        `json:"order_invoice_number"`
	OrderAmount        float64       `json:"order_amount"`
	Currency           string        `json:"currency"`
	OrderDescription   string        `json:"order_description"`
	OrderTaxAmount     *float64      `json:"order_tax_amount,omitempty"`
	CustomerID         *string       `json:"customer_id,omitempty"`
	SuccessURL         *string       `json:"success_url,omitempty"`
	ErrorURL           *string       `json:"error_url,omitempty"`
	CancelURL          *string       `json:"cancel_url,omitempty"`
	CustomData         *string       `json:"custom_data,omitempty"`
	Signature          string        `json:"signature"`
}

// FormValues returns the checkout fields as a map suitable for building
// an HTML form or URL-encoded POST body.
func (f *SignedCheckoutFields) FormValues() map[string]string {
	m := map[string]string{
		"merchant":             f.Merchant,
		"operation":            string(f.Operation),
		"order_invoice_number": f.OrderInvoiceNumber,
		"order_amount":         formatFloat(f.OrderAmount),
		"currency":             f.Currency,
		"order_description":    f.OrderDescription,
		"signature":            f.Signature,
	}
	if f.PaymentMethod != "" {
		m["payment_method"] = string(f.PaymentMethod)
	}
	if f.OrderTaxAmount != nil {
		m["order_tax_amount"] = formatFloat(*f.OrderTaxAmount)
	}
	if f.CustomerID != nil {
		m["customer_id"] = *f.CustomerID
	}
	if f.SuccessURL != nil {
		m["success_url"] = *f.SuccessURL
	}
	if f.ErrorURL != nil {
		m["error_url"] = *f.ErrorURL
	}
	if f.CancelURL != nil {
		m["cancel_url"] = *f.CancelURL
	}
	if f.CustomData != nil {
		m["custom_data"] = *f.CustomData
	}
	return m
}

// CheckoutService provides access to checkout URL generation and field signing.
type CheckoutService struct {
	client *Client
}

// InitCheckoutURL returns the URL for initiating a checkout session.
func (s *CheckoutService) InitCheckoutURL() string {
	return s.client.baseCheckoutURL + "/init"
}

// InitOneTimePaymentFields signs the given one-time payment fields and returns
// the complete set of fields including the HMAC-SHA256 signature.
func (s *CheckoutService) InitOneTimePaymentFields(fields OnetimePaymentFields) *SignedCheckoutFields {
	operation := fields.Operation
	if operation == "" {
		operation = OperationPurchase
	}

	signed := &SignedCheckoutFields{
		Merchant:           s.client.config.MerchantID,
		Operation:          operation,
		PaymentMethod:      fields.PaymentMethod,
		OrderInvoiceNumber: fields.OrderInvoiceNumber,
		OrderAmount:        fields.OrderAmount,
		Currency:           fields.Currency,
		OrderDescription:   fields.OrderDescription,
		OrderTaxAmount:     fields.OrderTaxAmount,
		CustomerID:         fields.CustomerID,
		SuccessURL:         fields.SuccessURL,
		ErrorURL:           fields.ErrorURL,
		CancelURL:          fields.CancelURL,
		CustomData:         fields.CustomData,
	}

	// Build map of signable fields.
	signableFields := map[string]string{
		"merchant":             signed.Merchant,
		"operation":            string(signed.Operation),
		"order_amount":         formatFloat(signed.OrderAmount),
		"currency":             signed.Currency,
		"order_invoice_number": signed.OrderInvoiceNumber,
		"order_description":    signed.OrderDescription,
	}
	if signed.PaymentMethod != "" {
		signableFields["payment_method"] = string(signed.PaymentMethod)
	}
	if signed.CustomerID != nil {
		signableFields["customer_id"] = *signed.CustomerID
	}
	if signed.SuccessURL != nil {
		signableFields["success_url"] = *signed.SuccessURL
	}
	if signed.ErrorURL != nil {
		signableFields["error_url"] = *signed.ErrorURL
	}
	if signed.CancelURL != nil {
		signableFields["cancel_url"] = *signed.CancelURL
	}

	signed.Signature = signFields(signableFields, s.client.config.SecretKey)
	return signed
}

// signFieldOrder defines the canonical order for signature computation.
var signFieldOrder = []string{
	"merchant",
	"env",
	"operation",
	"payment_method",
	"order_amount",
	"currency",
	"order_invoice_number",
	"order_description",
	"customer_id",
	"agreement_id",
	"agreement_name",
	"agreement_type",
	"agreement_payment_frequency",
	"agreement_amount_per_payment",
	"success_url",
	"error_url",
	"cancel_url",
	"order_id",
}

// signFields computes the HMAC-SHA256 signature for the given fields using the
// canonical field ordering, and returns the base64-encoded result.
func signFields(fields map[string]string, secretKey string) string {
	var parts []string
	for _, key := range signFieldOrder {
		val, ok := fields[key]
		if !ok {
			continue
		}
		parts = append(parts, key+"="+val)
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(strings.Join(parts, ",")))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// formatFloat formats a float64 without trailing zeros, matching JavaScript's
// Number.toString() behavior (e.g. 10000 → "10000", 100.5 → "100.5").
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
