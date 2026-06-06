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

	addLine("")

	addLine("TRACING:")
	addLine(fmt.Sprintf("DNS lookup:        %v", data.Trace.DNSDone.Sub(data.Trace.DNSStart)))
	addLine(fmt.Sprintf("TCP connect:       %v", data.Trace.ConnectDone.Sub(data.Trace.ConnectStart)))
	addLine(fmt.Sprintf("TLS handshake:     %v", data.Trace.TLSDone.Sub(data.Trace.TLSStart)))
	addLine(fmt.Sprintf("Server processing: %v", data.Trace.GotFirstResponseByte.Sub(data.Trace.GotConn)))
	addLine(fmt.Sprintf("Content transfer:  %v", data.Trace.Done.Sub(data.Trace.GotFirstResponseByte)))
	addLine(fmt.Sprintf("Total:             %v", data.Trace.Done.Sub(data.Trace.Start)))

	addLine("")

	addLine(fmt.Sprintf("Status: %s %s", data.Proto, data.Status))

	addLine("")

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

	addLine("")

	if data.Body != "" {
		addLine("BODY:")
		addLine(data.Body)
	}

	return strings.Join(lines, "\n")
}
