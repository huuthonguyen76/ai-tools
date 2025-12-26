package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type CustomHTTPClient struct {
	HTTPClient *http.Client
}

func NewHTTPClient(httpClient *http.Client) CustomHTTPClient {
	return CustomHTTPClient{
		HTTPClient: httpClient,
	}
}

func (s *CustomHTTPClient) Call(ctx context.Context, method, url string, requestBody *[]byte) (*http.Response, *[]byte, error) {
	var req *http.Request
	var err error
	if requestBody == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(*requestBody))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
		}
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close() // Without closing the leak will happened. When reading the io.ReadAll/

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, &body, nil
}
