package printer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

type jsonFormat struct {
	MetaInfo         metaInfo          `json:"meta_info"`
	Headers          map[string]string `json:"headers"`
	TransferEncoding []string          `json:"transer_encoding,omitempty"`
	Body             string            `json:"body,omitempty"`
}

type JsonFormatter struct {
}

func NewJSONFormatter() JsonFormatter {
	return JsonFormatter{}
}

func (j JsonFormatter) Format(
	description string,
	delta time.Duration,
	total time.Duration,
	proto string,
	status string,
	headers http.Header,
	transferEncoding []string,
	body string,
) string {
	output := jsonFormat{
		MetaInfo: metaInfo{
			Description: description,
			TimeDelta:   duration(delta),
			TimeTotal:   duration(total),
			Proto:       proto,
			Status:      status,
		},
		Headers:          convertHeaders(headers),
		TransferEncoding: transferEncoding,
		Body:             body,
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
