package request

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

type Interceptor interface {
	Before()
	After(rsp *http.Response)
	Redirect(rsp *http.Response)
}

type Request struct {
	interceptor Interceptor
}

type RedirectFunc func(request *http.Request, via []*http.Request) error

func New(interceptor Interceptor) *Request {
	return &Request{
		interceptor: interceptor,
	}
}

func (r *Request) Do(ctx context.Context, method, url string, opts ...Option) error {
	var reqOptions options

	for _, opt := range opts {
		opt(&reqOptions)
	}

	client := &http.Client{
		CheckRedirect: r.doRedirect(),
	}

	if reqOptions.ForceHTTP11 {
		client.Transport = &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		}
	}

	return r.doImpl(ctx, method, url, client)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func (r *Request) doImpl(ctx context.Context, method, url string, doer doer) error {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("create new request: %w", err)
	}

	r.interceptor.Before()

	response, err := doer.Do(request)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}

	defer response.Body.Close()

	r.interceptor.After(response)

	return nil
}

func (r *Request) doRedirect() RedirectFunc {
	return func(request *http.Request, via []*http.Request) error {
		r.interceptor.Redirect(request.Response)

		return nil
	}
}
