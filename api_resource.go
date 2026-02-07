package sepay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type apiResource struct {
	client *Client
}

func (a *apiResource) authHeader() string {
	creds := a.client.config.MerchantID + ":" + a.client.config.SecretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func (a *apiResource) doRequest(ctx context.Context, method, endpoint string, query url.Values, body io.Reader) (*Response, error) {
	rawURL := a.client.baseAPIURL + "/" + endpoint
	if len(query) > 0 {
		rawURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("sepay: creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", a.authHeader())

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sepay: executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sepay: reading response body: %w", err)
	}

	r := &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}

	if resp.StatusCode >= 400 {
		return r, &APIError{
			StatusCode: resp.StatusCode,
			Body:       respBody,
		}
	}

	return r, nil
}

func (a *apiResource) doRequestJSON(ctx context.Context, method, endpoint string, body any) (*Response, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("sepay: marshaling request body: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	return a.doRequest(ctx, method, endpoint, nil, reader)
}
