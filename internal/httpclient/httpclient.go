package httpclient

import (
	"bytes"
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
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return h.c.Post(url, ApplicationJSON, bytes.NewBuffer(jsonBody))
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
