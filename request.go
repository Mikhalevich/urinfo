package main

import (
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

func (r *Request) doRedirect() RedirectFunc {
	return func(request *http.Request, via []*http.Request) error {
		currentTime := time.Now()
		r.P.Print("redirect", currentTime.Sub(r.PreviousTime), currentTime.Sub(r.StartTime), request.Response)
		r.PreviousTime = currentTime

		return nil
	}
}

func (r *Request) Do(method, url string) error {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		CheckRedirect: r.doRedirect(),
	}

	r.StartTime = time.Now()
	r.PreviousTime = r.StartTime

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	currentTime := time.Now()
	r.P.Print("result", currentTime.Sub(r.PreviousTime), currentTime.Sub(r.StartTime), response)

	return nil
}
