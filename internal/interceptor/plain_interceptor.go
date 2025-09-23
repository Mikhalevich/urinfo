package interceptor

import (
	"io"
	"log"
	"net/http"
	"time"
)

type PlainInterceptor struct {
	isPrintBody  bool
	startTime    time.Time
	previousTime time.Time
}

func NewPlainInterceptor(isPrintBody bool) *PlainInterceptor {
	return &PlainInterceptor{
		isPrintBody: isPrintBody,
	}
}

func (p *PlainInterceptor) Before() {
	p.startTime = time.Now()
	p.previousTime = p.startTime
}

func (p *PlainInterceptor) After(rsp *http.Response) {
	now := time.Now()

	p.print("result", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)
}

func (p *PlainInterceptor) Redirect(rsp *http.Response) {
	now := time.Now()

	p.print("redirect", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)

	p.previousTime = now
}

func (p *PlainInterceptor) print(
	description string,
	delta time.Duration,
	total time.Duration,
	response *http.Response,
) {
	log.Printf("<<<<<<<<<<<<<<<<<<<<<<<< %s delta = %s total = %s\n", description, delta, total)

	log.Printf("Status: %s %s\n", response.Proto, response.Status)

	log.Println("HEADERS:")

	for key, value := range response.Header {
		log.Printf("%s => %s\n", key, value)
	}

	if response.TransferEncoding != nil {
		log.Println("Transfer Encoding:")

		for _, v := range response.TransferEncoding {
			log.Println(v)
		}
	}

	if p.isPrintBody && response.Body != nil {
		log.Println("BODY:")

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)

			return
		}

		log.Println(string(body))
	}
}
