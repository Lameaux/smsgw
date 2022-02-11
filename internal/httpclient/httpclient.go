package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

const (
	ApplicationJSON = "application/json"
)

type HTTPClient struct {
	c *http.Client
}

type Builder struct {
	connectionTimeout time.Duration
	tlsTimeout        time.Duration
	readTimeout       time.Duration
}

func (h *HTTPClient) Post(url string, body interface{}) (*http.Response, error) {
	return h.SendWithContext(context.Background(), http.MethodPost, url, body)
}

func (h *HTTPClient) SendWithContext(ctx context.Context, method string, url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ApplicationJSON)

	return h.c.Do(req)
}

func newHTTPClient(b *Builder) *HTTPClient {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: b.connectionTimeout,
		}).Dial,
		TLSHandshakeTimeout: b.tlsTimeout,
	}

	c := &http.Client{
		Timeout:   b.readTimeout,
		Transport: netTransport,
	}

	return &HTTPClient{c}
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ConnectionTimeout(t time.Duration) *Builder {
	b.connectionTimeout = t

	return b
}

func (b *Builder) TLSTimeout(t time.Duration) *Builder {
	b.tlsTimeout = t

	return b
}

func (b *Builder) ReadTimeout(t time.Duration) *Builder {
	b.readTimeout = t

	return b
}

func (b *Builder) Build() *HTTPClient {
	return newHTTPClient(b)
}
