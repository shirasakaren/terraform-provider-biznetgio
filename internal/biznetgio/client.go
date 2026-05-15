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
