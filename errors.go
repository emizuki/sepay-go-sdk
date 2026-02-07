package sepay

import "fmt"

// ConfigError is returned when the client configuration is invalid.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("sepay: config error: %s: %s", e.Field, e.Message)
}

// APIError is returned when the API returns an HTTP status code >= 400.
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("sepay: api error: status %d: %s", e.StatusCode, string(e.Body))
}
