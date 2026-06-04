package formatter

import (
	"fmt"
	"strings"

	"github.com/Mikhalevich/urinfo/internal/interceptor/printer"
)

type PlainFormatter struct {
}

func NewPlainFormatter() PlainFormatter {
	return PlainFormatter{}
}

func (p PlainFormatter) Format(data printer.ResponseData) string {
	var (
		lines   []string
		addLine = func(line string) {
			lines = append(lines, line)
		}
	)

	addLine(fmt.Sprintf("<<<<<<<<<<<<<<<<<<<<<<<< %s delta = %s total = %s",
		data.Description, data.Delta, data.Total))

	addLine(fmt.Sprintf("Status: %s %s", data.Proto, data.Status))

	addLine("HEADERS:")

	for key, value := range data.Headers {
		addLine(fmt.Sprintf("%s => %s", key, value))
	}

	if len(data.TransferEncoding) != 0 {
		addLine("Transfer Encoding:")

		for _, v := range data.TransferEncoding {
			addLine(v)
		}
	}

	if data.Body != "" {
		addLine("BODY:")
		addLine(data.Body)
	}

	return strings.Join(lines, "\n")
}
