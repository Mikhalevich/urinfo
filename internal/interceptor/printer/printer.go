package printer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Formatter interface {
	Format(
		description string,
		delta time.Duration,
		total time.Duration,
		proto string,
		status string,
		headers http.Header,
		transferEncoding []string,
		body string,
	) string
}

type Printer struct {
	isPrintBody  bool
	formatter    Formatter
	startTime    time.Time
	previousTime time.Time
}

func NewPrinter(isPrintBody bool, formatter Formatter) *Printer {
	return &Printer{
		isPrintBody: isPrintBody,
		formatter:   formatter,
	}
}

func NewPlainPrinter(isPrintBody bool) *Printer {
	return NewPrinter(isPrintBody, NewPlainFormatter())
}

func (p *Printer) Before() {
	p.startTime = time.Now()
	p.previousTime = p.startTime
}

func (p *Printer) After(rsp *http.Response) {
	now := time.Now()

	p.print("result", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)
}

func (p *Printer) Redirect(rsp *http.Response) {
	now := time.Now()

	p.print("redirect", now.Sub(p.previousTime), now.Sub(p.startTime), rsp)

	p.previousTime = now
}

func (p *Printer) print(
	description string,
	delta time.Duration,
	total time.Duration,
	rsp *http.Response,
) {
	body, err := p.responseBody(rsp)
	if err != nil {
		log.Println(err)
	}

	output := p.formatter.Format(
		description,
		delta,
		total,
		rsp.Proto,
		rsp.Status,
		rsp.Header,
		rsp.TransferEncoding,
		body,
	)

	log.Println(output)
}

func (p *Printer) responseBody(rsp *http.Response) (string, error) {
	if !p.isPrintBody || rsp.Body == nil {
		return "", nil
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return string(body), nil
}
