package sepay

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("valid config with defaults", func(t *testing.T) {
		c, err := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.config.APIVersion != APIVersionV1 {
			t.Errorf("expected APIVersion %q, got %q", APIVersionV1, c.config.APIVersion)
		}
		if c.config.CheckoutVersion != CheckoutVersionV1 {
			t.Errorf("expected CheckoutVersion %q, got %q", CheckoutVersionV1, c.config.CheckoutVersion)
		}
		if c.Order == nil {
			t.Error("expected Order service to be initialized")
		}
		if c.Checkout == nil {
			t.Error("expected Checkout service to be initialized")
		}
	})

	t.Run("empty merchant ID", func(t *testing.T) {
		_, err := NewClient(Config{
			Env:       Sandbox,
			SecretKey: "secret456",
		})
		if err == nil {
			t.Fatal("expected error for empty MerchantID")
		}
		cfgErr, ok := err.(*ConfigError)
		if !ok {
			t.Fatalf("expected *ConfigError, got %T", err)
		}
		if cfgErr.Field != "MerchantID" {
			t.Errorf("expected field %q, got %q", "MerchantID", cfgErr.Field)
		}
	})

	t.Run("empty secret key", func(t *testing.T) {
		_, err := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "merchant123",
		})
		if err == nil {
			t.Fatal("expected error for empty SecretKey")
		}
		cfgErr, ok := err.(*ConfigError)
		if !ok {
			t.Fatalf("expected *ConfigError, got %T", err)
		}
		if cfgErr.Field != "SecretKey" {
			t.Errorf("expected field %q, got %q", "SecretKey", cfgErr.Field)
		}
	})

	t.Run("unsupported API version", func(t *testing.T) {
		_, err := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
			APIVersion: "v2",
		})
		if err == nil {
			t.Fatal("expected error for unsupported APIVersion")
		}
		cfgErr, ok := err.(*ConfigError)
		if !ok {
			t.Fatalf("expected *ConfigError, got %T", err)
		}
		if cfgErr.Field != "APIVersion" {
			t.Errorf("expected field %q, got %q", "APIVersion", cfgErr.Field)
		}
	})

	t.Run("unsupported checkout version", func(t *testing.T) {
		_, err := NewClient(Config{
			Env:             Sandbox,
			MerchantID:      "merchant123",
			SecretKey:       "secret456",
			CheckoutVersion: "v2",
		})
		if err == nil {
			t.Fatal("expected error for unsupported CheckoutVersion")
		}
		cfgErr, ok := err.(*ConfigError)
		if !ok {
			t.Fatalf("expected *ConfigError, got %T", err)
		}
		if cfgErr.Field != "CheckoutVersion" {
			t.Errorf("expected field %q, got %q", "CheckoutVersion", cfgErr.Field)
		}
	})

	t.Run("sandbox base URLs", func(t *testing.T) {
		c, err := NewClient(Config{
			Env:        Sandbox,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedAPI := "https://pgapi-sandbox.sepay.vn/v1"
		if c.baseAPIURL != expectedAPI {
			t.Errorf("expected baseAPIURL %q, got %q", expectedAPI, c.baseAPIURL)
		}
		expectedCheckout := "https://pay-sandbox.sepay.vn/v1/checkout"
		if c.baseCheckoutURL != expectedCheckout {
			t.Errorf("expected baseCheckoutURL %q, got %q", expectedCheckout, c.baseCheckoutURL)
		}
	})

	t.Run("production base URLs", func(t *testing.T) {
		c, err := NewClient(Config{
			Env:        Production,
			MerchantID: "merchant123",
			SecretKey:  "secret456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedAPI := "https://pgapi.sepay.vn/v1"
		if c.baseAPIURL != expectedAPI {
			t.Errorf("expected baseAPIURL %q, got %q", expectedAPI, c.baseAPIURL)
		}
		expectedCheckout := "https://pay.sepay.vn/v1/checkout"
		if c.baseCheckoutURL != expectedCheckout {
			t.Errorf("expected baseCheckoutURL %q, got %q", expectedCheckout, c.baseCheckoutURL)
		}
	})
}
