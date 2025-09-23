package printer

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PlainFormatter struct {
}

func NewPlainFormatter() PlainFormatter {
	return PlainFormatter{}
}

func (p PlainFormatter) Format(
	description string,
	delta time.Duration,
	total time.Duration,
	proto string,
	status string,
	headers http.Header,
	transferEncoding []string,
	body string,
) string {
	var (
		lines   []string
		addLine = func(line string) {
			lines = append(lines, line)
		}
	)

	addLine(fmt.Sprintf("<<<<<<<<<<<<<<<<<<<<<<<< %s delta = %s total = %s", description, delta, total))

	addLine(fmt.Sprintf("Status: %s %s", proto, status))

	addLine("HEADERS:")

	for key, value := range headers {
		addLine(fmt.Sprintf("%s => %s", key, value))
	}

	if transferEncoding != nil {
		addLine("Transfer Encoding:")

		for _, v := range transferEncoding {
			addLine(v)
		}
	}

	if body != "" {
		addLine("BODY:")
		addLine(body)
	}

	return strings.Join(lines, "\n")
}
