package plain

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Plain struct {
	isPrintBody  bool
	startTime    time.Time
	previousTime time.Time
}

func New(isPrintBody bool) *Plain {
	return &Plain{
		isPrintBody: isPrintBody,
	}
}

func (p *Plain) Before() {
	p.startTime = time.Now()
	p.previousTime = p.startTime
}

func (p *Plain) After(rsp *http.Response) {
	now := time.Now()

	p.print("result", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)
}

func (p *Plain) Redirect(rsp *http.Response) {
	now := time.Now()

	p.print("redirect", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)

	p.previousTime = now
}

func (p *Plain) print(
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
