package formatter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Mikhalevich/urinfo/internal/interceptor/printer"
)

type duration time.Duration

func (d duration) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(time.Duration(d).String())
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return b, nil
}

type metaInfo struct {
	Description string   `json:"description"`
	TimeDelta   duration `json:"time_delta"`
	TimeTotal   duration `json:"time_total"`
	Proto       string   `json:"proto"`
	Status      string   `json:"status"`
}

type tracing struct {
	DNSLookup        duration `json:"dns_lookup"`
	TCPConnect       duration `json:"tcp_connect"`
	TLSHandshake     duration `json:"tls_handshake"`
	ServerProcessing duration `json:"server_processing"`
	ContentTransfer  duration `json:"content_transfer"`
	Total            duration `json:"total"`
}

type jsonFormat struct {
	MetaInfo         metaInfo          `json:"meta_info"`
	Tracing          tracing           `json:"tracing"`
	Headers          map[string]string `json:"headers"`
	TransferEncoding []string          `json:"transer_encoding,omitempty"`
	Body             string            `json:"body,omitempty"`
}

type JsonFormatter struct {
}

func NewJSONFormatter() JsonFormatter {
	return JsonFormatter{}
}

func (j JsonFormatter) Format(data printer.ResponseData) string {
	output := jsonFormat{
		MetaInfo: metaInfo{
			Description: data.Description,
			TimeDelta:   duration(data.Delta),
			TimeTotal:   duration(data.Total),
			Proto:       data.Proto,
			Status:      data.Status,
		},
		Tracing: tracing{
			DNSLookup:        duration(data.Trace.DNSDone.Sub(data.Trace.DNSStart)),
			TCPConnect:       duration(data.Trace.ConnectDone.Sub(data.Trace.ConnectStart)),
			TLSHandshake:     duration(data.Trace.TLSDone.Sub(data.Trace.TLSStart)),
			ServerProcessing: duration(data.Trace.GotFirstResponseByte.Sub(data.Trace.GotConn)),
			ContentTransfer:  duration(data.Trace.Done.Sub(data.Trace.GotFirstResponseByte)),
			Total:            duration(data.Trace.Done.Sub(data.Trace.Start)),
		},
		Headers:          convertHeaders(data.Headers),
		TransferEncoding: data.TransferEncoding,
		Body:             data.Body,
	}

	b, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return ""
	}

	return string(b)
}

func convertHeaders(headers http.Header) map[string]string {
	if len(headers) == 0 {
		return nil
	}

	converted := make(map[string]string, len(headers))

	for k, v := range headers {
		converted[k] = strings.Join(v, "")
	}

	return converted
}
