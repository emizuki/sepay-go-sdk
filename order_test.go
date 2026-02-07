package sepay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	c, err := NewClient(Config{
		Env:        Sandbox,
		MerchantID: "merchant123",
		SecretKey:  "secret456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.baseAPIURL = ts.URL
	return c, ts
}

func TestOrderService_All(t *testing.T) {
	t.Run("basic request", func(t *testing.T) {
		c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/order" {
				t.Errorf("expected path /order, got %s", r.URL.Path)
			}

			// Verify auth header.
			expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("merchant123:secret456"))
			if r.Header.Get("Authorization") != expectedAuth {
				t.Errorf("expected auth %q, got %q", expectedAuth, r.Header.Get("Authorization"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"data":[]}`))
		})
		defer ts.Close()

		resp, err := c.Order.All(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("with query params", func(t *testing.T) {
		c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("per_page") != "10" {
				t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
			}
			if q.Get("order_status") != "COMPLETED" {
				t.Errorf("expected order_status=COMPLETED, got %q", q.Get("order_status"))
			}
			if q.Get("sort[created_at]") != "desc" {
				t.Errorf("expected sort[created_at]=desc, got %q", q.Get("sort[created_at]"))
			}
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		})
		defer ts.Close()

		_, err := c.Order.All(context.Background(), &OrderQueryParams{
			PerPage:       Int(10),
			OrderStatus:   String("COMPLETED"),
			SortCreatedAt: String("desc"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestOrderService_Retrieve(t *testing.T) {
	c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/order/detail/INV-001" {
			t.Errorf("expected path /order/detail/INV-001, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"order_invoice_number":"INV-001"}`))
	})
	defer ts.Close()

	resp, err := c.Order.Retrieve(context.Background(), "INV-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestOrderService_VoidTransaction(t *testing.T) {
	c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/order/voidTransaction" {
			t.Errorf("expected path /order/voidTransaction, got %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["order_invoice_number"] != "INV-001" {
			t.Errorf("expected order_invoice_number INV-001, got %q", body["order_invoice_number"])
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"voided"}`))
	})
	defer ts.Close()

	resp, err := c.Order.VoidTransaction(context.Background(), "INV-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestOrderService_Cancel(t *testing.T) {
	c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/order/cancel" {
			t.Errorf("expected path /order/cancel, got %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["order_invoice_number"] != "INV-002" {
			t.Errorf("expected order_invoice_number INV-002, got %q", body["order_invoice_number"])
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"cancelled"}`))
	})
	defer ts.Close()

	resp, err := c.Order.Cancel(context.Background(), "INV-002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestOrderService_APIError(t *testing.T) {
	c, ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"unauthorized"}`))
	})
	defer ts.Close()

	resp, err := c.Order.All(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}

	// Response should still be available.
	if resp == nil {
		t.Fatal("expected non-nil response even on error")
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected response status 401, got %d", resp.StatusCode)
	}
}
