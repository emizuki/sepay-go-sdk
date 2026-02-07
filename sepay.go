package sepay

import "net/http"

// Environment represents the SePay environment.
type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

// APIVersion represents the SePay API version.
type APIVersion string

const (
	APIVersionV1 APIVersion = "v1"
)

// CheckoutVersion represents the SePay checkout version.
type CheckoutVersion string

const (
	CheckoutVersionV1 CheckoutVersion = "v1"
)

// Config holds the configuration for a SePay client.
type Config struct {
	Env             Environment
	MerchantID      string
	SecretKey       string
	APIVersion      APIVersion
	CheckoutVersion CheckoutVersion
}

// Client is the SePay payment gateway client.
type Client struct {
	Order    *OrderService
	Checkout *CheckoutService

	config          Config
	baseAPIURL      string
	baseCheckoutURL string
	httpClient      *http.Client
}

// NewClient creates a new SePay client with the given configuration.
func NewClient(cfg Config) (*Client, error) {
	if cfg.MerchantID == "" {
		return nil, &ConfigError{Field: "MerchantID", Message: "must not be empty"}
	}
	if cfg.SecretKey == "" {
		return nil, &ConfigError{Field: "SecretKey", Message: "must not be empty"}
	}

	if cfg.APIVersion == "" {
		cfg.APIVersion = APIVersionV1
	}
	if cfg.CheckoutVersion == "" {
		cfg.CheckoutVersion = CheckoutVersionV1
	}

	if cfg.APIVersion != APIVersionV1 {
		return nil, &ConfigError{Field: "APIVersion", Message: "unsupported version"}
	}
	if cfg.CheckoutVersion != CheckoutVersionV1 {
		return nil, &ConfigError{Field: "CheckoutVersion", Message: "unsupported version"}
	}

	var baseAPIURL, baseCheckoutURL string
	if cfg.Env == Sandbox {
		baseAPIURL = "https://pgapi-sandbox.sepay.vn/" + string(cfg.APIVersion)
		baseCheckoutURL = "https://pay-sandbox.sepay.vn/" + string(cfg.CheckoutVersion) + "/checkout"
	} else {
		baseAPIURL = "https://pgapi.sepay.vn/" + string(cfg.APIVersion)
		baseCheckoutURL = "https://pay.sepay.vn/" + string(cfg.CheckoutVersion) + "/checkout"
	}

	c := &Client{
		config:          cfg,
		baseAPIURL:      baseAPIURL,
		baseCheckoutURL: baseCheckoutURL,
		httpClient:      &http.Client{},
	}

	c.Order = &OrderService{api: apiResource{client: c}}
	c.Checkout = &CheckoutService{client: c}

	return c, nil
}

// SetHTTPClient sets a custom HTTP client for the SePay client.
// This is useful for testing or for using custom transports.
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}
