package sepay

import (
	"encoding/json"
	"net/http"
)

// Response wraps an HTTP response from the SePay API.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// DecodeJSON decodes the response body into the given value.
func (r *Response) DecodeJSON(v any) error {
	return json.Unmarshal(r.Body, v)
}
