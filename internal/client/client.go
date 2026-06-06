package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Mikhalevich/urinfo/internal/trace"
)

type Interceptor interface {
	Before()
	After(rsp *http.Response, tracing trace.Trace)
	Redirect(rsp *http.Response, tracing trace.Trace)
}

type RedirectFunc func(request *http.Request, via []*http.Request) error

type Client struct {
	intercetor Interceptor
	transport  *Transport
	httpClient *http.Client
}

func New(intercetor Interceptor, opts ...Option) *Client {
	var clientOptions options

	for _, opt := range opts {
		opt(&clientOptions)
	}

	httpTransport := http.DefaultTransport

	if clientOptions.ForceHTTP11 {
		httpTransport = &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		}
	}

	transport := NewTransport(httpTransport, intercetor)

	return &Client{
		intercetor: intercetor,
		transport:  transport,
		httpClient: &http.Client{
			Transport:     transport,
			CheckRedirect: doRedirect(intercetor, transport),
		},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.intercetor.Before()

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}

	c.intercetor.After(rsp, *c.transport.tracing)

	return rsp, nil
}

func doRedirect(interceptor Interceptor, transport *Transport) RedirectFunc {
	return func(request *http.Request, via []*http.Request) error {
		interceptor.Redirect(request.Response, *transport.tracing)

		return nil
	}
}

func Do(ctx context.Context, interceptor Interceptor, method, url string, opts ...Option) error {
	client := New(interceptor, opts...)

	rsp, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("create new request: %w", err)
	}

	response, err := client.Do(rsp)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}

	defer response.Body.Close()

	return nil
}
