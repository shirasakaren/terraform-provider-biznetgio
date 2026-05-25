package biznetgio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.portal.biznetgio.com/v1"
const defaultTimeout = 30 * time.Second

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, apiKey string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   apiKey,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) Neolite() *NeoliteService {
	return &NeoliteService{client: c}
}

func (c *Client) NeolitePro() *NeoliteProService {
	return &NeoliteProService{client: c}
}

func (c *Client) Baremetal() *BaremetalService {
	return &BaremetalService{client: c}
}

func (c *Client) BaremetalAdditionalIP() *BaremetalAdditionalIPService {
	return &BaremetalAdditionalIPService{client: c}
}

func (c *Client) BaremetalElasticStorage() *BaremetalElasticStorageService {
	return &BaremetalElasticStorageService{client: c}
}

func (c *Client) GPU() *GPUService {
	return &GPUService{client: c}
}

func (c *Client) ObjectStorage() *ObjectStorageService {
	return &ObjectStorageService{client: c}
}

func (c *Client) doJSON(ctx context.Context, method, path string, in any, out any) error {
	var body []byte
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = b
	}
	return c.do(ctx, method, path, body, "application/json", out)
}

func (c *Client) do(ctx context.Context, method, path string, body []byte, contentType string, out any) error {
	_, err := c.send(ctx, method, path, body, contentType, func(status int, raw []byte) (bool, error) {
		if status >= 200 && status <= 299 {
			return true, c.unwrap(status, raw, out)
		}
		return false, newAPIError(status, raw)
	})
	return err
}

func (c *Client) raw(ctx context.Context, method, path string, body []byte, contentType string) ([]byte, error) {
	return c.send(ctx, method, path, body, contentType, func(status int, raw []byte) (bool, error) {
		if status >= 200 && status <= 299 {
			return true, nil
		}
		return false, newAPIError(status, raw)
	})
}

func (c *Client) send(ctx context.Context, method, path string, body []byte, contentType string, handle func(status int, raw []byte) (bool, error)) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= retryAttempts; attempt++ {
		if attempt > 0 {
			if err := sleepWithCtx(ctx, retryBackoff(attempt)); err != nil {
				return nil, err
			}
		}
		status, raw, err := c.execute(ctx, method, path, body, contentType)
		if err != nil {
			return nil, err
		}
		done, err := handle(status, raw)
		if done {
			return raw, err
		}
		lastErr = err
		if !retryable(err) {
// wip 341
