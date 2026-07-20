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
// wip 849
// wip 946
// wip 962
