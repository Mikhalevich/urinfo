package printer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Mikhalevich/urinfo/internal/trace"
)

type ResponseData struct {
	Description      string
	Delta            time.Duration
	Total            time.Duration
	Proto            string
	Status           string
	Headers          http.Header
	TransferEncoding []string
	Body             string
	Trace            trace.Trace
}

type Formatter interface {
	Format(data ResponseData) string
}

type Printer struct {
	isPrintBody   bool
	formatter     Formatter
	startTime     time.Time
	previousTime  time.Time
	responseSteps []ResponseData
}

func NewPrinter(isPrintBody bool, formatter Formatter) *Printer {
	return &Printer{
		isPrintBody: isPrintBody,
		formatter:   formatter,
	}
}

func (p *Printer) Before() {
	if p.responseSteps != nil {
		p.responseSteps = p.responseSteps[0:]
	}

	p.startTime = time.Now()
	p.previousTime = p.startTime
}

func (p *Printer) After(rsp *http.Response, tracing trace.Trace) {
	now := time.Now()

	p.addResponseStep("result", now.Sub(p.previousTime), now.Sub(p.startTime), rsp, tracing)

	p.printSteps()
}

func (p *Printer) Redirect(rsp *http.Response, tracing trace.Trace) {
	now := time.Now()

	p.addResponseStep("redirect", now.Sub(p.previousTime), now.Sub(p.startTime), rsp, tracing)

	p.previousTime = now
}

func (p *Printer) addResponseStep(
	description string,
	delta time.Duration,
	total time.Duration,
	rsp *http.Response,
	tracing trace.Trace,
) {
	body, err := p.responseBody(rsp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	p.responseSteps = append(p.responseSteps, ResponseData{
		Description:      description,
		Delta:            delta,
		Total:            total,
		Proto:            rsp.Proto,
		Status:           rsp.Status,
		Headers:          rsp.Header,
		TransferEncoding: rsp.TransferEncoding,
		Body:             body,
		Trace:            tracing,
	})
}

func (p *Printer) printSteps() {
	for _, step := range p.responseSteps {
		fmt.Fprintln(os.Stdout, p.formatter.Format(step))
	}
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
