package client

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/Mikhalevich/urinfo/internal/trace"
)

type Transport struct {
	base    http.RoundTripper
	tracing *trace.Trace
}

func NewTransport(base http.RoundTripper, interceptor Interceptor) *Transport {
	return &Transport{
		base: base,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		tracing, clientTrace = trace.New()
		ctx                  = httptrace.WithClientTrace(req.Context(), clientTrace)
	)

	t.tracing = tracing
	tracing.Start = time.Now()

	rsp, err := t.base.RoundTrip(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("round trip: %w", err)
	}

	tracing.Done = time.Now()

	return rsp, nil
}
