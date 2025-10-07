package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/Mikhalevich/urinfo/internal/interceptor/printer"
	"github.com/Mikhalevich/urinfo/internal/request"
)

type OutputFormat int

const (
	PlainFormat OutputFormat = iota + 1
	JSONFormat
)

type Params struct {
	URL         string
	Method      string
	PrintBody   bool
	ForceHTTP11 bool
	Format      OutputFormat
}

func getURL() (string, error) {
	if flag.NArg() <= 0 {
		return "", errors.New("no url specified")
	}

	urlString := flag.Arg(0)

	uri, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	if uri.Scheme == "" {
		urlString = "http://" + urlString
	}

	return urlString, nil
}

func parseArguments() (*Params, error) {
	customMethod := flag.String("method", "", "custom method")
	isGet := flag.Bool("get", false, "get method")
	isPost := flag.Bool("post", false, "post method")
	isHead := flag.Bool("head", false, "head method")
	noBody := flag.Bool("nobody", false, "print result without body")
	http11 := flag.Bool("http11", false, "use HTTP/1.1 protocol")
	jsonFormat := flag.Bool("json", false, "json output")

	flag.Parse()

	urlString, err := getURL()
	if err != nil {
		return nil, err
	}

	method := "GET"

	switch {
	case *isGet:
		method = "GET"
	case *isPost:
		method = "POST"
	case *isHead:
		method = "HEAD"
	case *customMethod != "":
		method = *customMethod
	}

	format := PlainFormat

	if *jsonFormat {
		format = JSONFormat
	}

	return &Params{
		URL:         urlString,
		Method:      method,
		PrintBody:   !*noBody,
		ForceHTTP11: *http11,
		Format:      format,
	}, nil
}

func main() {
	params, err := parseArguments()
	if err != nil {
		log.Fatalln(err)

		return
	}

	r := request.New(makePrinter(params.Format, params.PrintBody))
	if err := r.Do(
		context.Background(),
		params.Method,
		params.URL,
		request.WithForceHTTP11(params.ForceHTTP11),
	); err != nil {
		log.Fatalln(err)

		return
	}
}

func makePrinter(format OutputFormat, isPrintBody bool) *printer.Printer {
	switch format {
	case PlainFormat:
		return printer.NewPlainPrinter(isPrintBody)

	case JSONFormat:
		return printer.NewJSONPrinter(isPrintBody)
	}

	return nil
}
