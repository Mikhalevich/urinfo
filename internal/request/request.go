package request

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/Mikhalevich/urinfo/internal/trace"
)

type Interceptor interface {
	Before()
	After(rsp *http.Response, tracing trace.Trace)
	Redirect(rsp *http.Response, tracing trace.Trace)
}

type Request struct {
	interceptor Interceptor
	tracing     *trace.Trace
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

	var (
		tracing, clientTrace = trace.New()
	)

	r.tracing = tracing

	return r.doImpl(httptrace.WithClientTrace(ctx, clientTrace), method, url, client)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func (r *Request) doImpl(
	ctx context.Context,
	method, url string,
	doer doer,
) error {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("create new request: %w", err)
	}

	r.interceptor.Before()
	r.tracing.Start = time.Now()

	response, err := doer.Do(request)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}

	defer response.Body.Close()

	r.tracing.Done = time.Now()

	r.interceptor.After(response, copyTrace(r.tracing))

	return nil
}

func (r *Request) doRedirect() RedirectFunc {
	return func(request *http.Request, via []*http.Request) error {
		r.tracing.Done = time.Now()

		r.interceptor.Redirect(request.Response, copyTrace(r.tracing))

		r.tracing.Start = time.Now()

		return nil
	}
}

func copyTrace(tracing *trace.Trace) trace.Trace {
	if tracing != nil {
		return *tracing
	}

	return trace.Trace{}
}
