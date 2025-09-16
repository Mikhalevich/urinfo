package consoleinterceptor

import (
	"io"
	"log"
	"net/http"
	"time"
)

type ConsoleInterceptor struct {
	isPrintBody  bool
	startTime    time.Time
	previousTime time.Time
}

func New(isPrintBody bool) *ConsoleInterceptor {
	return &ConsoleInterceptor{
		isPrintBody: isPrintBody,
	}
}

func (c *ConsoleInterceptor) Before() {
	c.startTime = time.Now()
	c.previousTime = c.startTime
}

func (c *ConsoleInterceptor) After(rsp *http.Response) {
	now := time.Now()

	c.print("result", now.Sub(c.previousTime), now.Sub(c.startTime), rsp)
}

func (c *ConsoleInterceptor) Redirect(rsp *http.Response) {
	now := time.Now()

	c.print("redirect", now.Sub(c.previousTime), now.Sub(c.startTime), rsp)

	c.previousTime = now
}

func (c *ConsoleInterceptor) print(
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

	if c.isPrintBody && response.Body != nil {
		log.Println("BODY:")

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)

			return
		}

		log.Println(string(body))
	}
}
