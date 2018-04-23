package main

import (
	"io"
	"net/http"
	"time"
)

type Request struct {
	StartTime    time.Time
	PreviousTime time.Time
	P            Printer
}

type RedirectFunc func(request *http.Request, via []*http.Request) error

func NewRequest(p Printer) *Request {
	if p == nil {
		p = &EmptyPrint{}
	}

	return &Request{
		P: p,
	}
}

func (r *Request) doRedirect(printBody bool) RedirectFunc {
	return func(request *http.Request, via []*http.Request) error {
		var body *io.ReadCloser = nil
		if printBody {
			body = &request.Response.Body
		}

		currentTime := time.Now()
		r.P.Print("redirect", currentTime.Sub(r.PreviousTime), currentTime.Sub(r.StartTime), request.Response.Status, &request.Response.Header, body)
		r.PreviousTime = currentTime

		return nil
	}
}

func (r *Request) Do(params *UriParams) error {
	request, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		CheckRedirect: r.doRedirect(params.PrintBody),
	}

	r.StartTime = time.Now()
	r.PreviousTime = r.StartTime

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var body *io.ReadCloser = nil
	if params.PrintBody {
		body = &response.Body
	}

	currentTime := time.Now()
	r.P.Print("result", currentTime.Sub(r.PreviousTime), currentTime.Sub(r.StartTime), response.Status, &response.Header, body)

	return nil
}
