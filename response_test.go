package sepay

import (
	"testing"
)

func TestResponse_DecodeJSON(t *testing.T) {
	t.Run("decode valid JSON", func(t *testing.T) {
		r := &Response{
			StatusCode: 200,
			Body:       []byte(`{"id":1,"name":"test"}`),
		}
		var result struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := r.DecodeJSON(&result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != 1 {
			t.Errorf("expected ID 1, got %d", result.ID)
		}
		if result.Name != "test" {
			t.Errorf("expected Name %q, got %q", "test", result.Name)
		}
	})

	t.Run("decode invalid JSON", func(t *testing.T) {
		r := &Response{
			StatusCode: 200,
			Body:       []byte(`not json`),
		}
		var result map[string]any
		if err := r.DecodeJSON(&result); err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("decode empty body", func(t *testing.T) {
		r := &Response{
			StatusCode: 200,
			Body:       []byte{},
		}
		var result map[string]any
		if err := r.DecodeJSON(&result); err == nil {
			t.Fatal("expected error for empty body")
		}
	})
}
