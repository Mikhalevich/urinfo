package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type UriParams struct {
	Url       string
	Method    string
	PrintBody bool
}

type Printer interface {
	Print(description string, delta time.Duration, total time.Duration, status string, headers *http.Header, body *io.ReadCloser)
}

func getUrl() (string, error) {
	if flag.NArg() <= 0 {
		return "", errors.New("No url specified")
	}

	urlString := flag.Arg(0)
	uri, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if uri.Scheme == "" {
		urlString = "http://" + urlString
	}

	return urlString, nil
}

func parseArguments() (*UriParams, error) {
	customMethod := flag.String("method", "", "custom method")
	isGet := flag.Bool("get", false, "get method")
	isPost := flag.Bool("post", false, "post method")
	isHead := flag.Bool("head", false, "head method")
	noBody := flag.Bool("nobody", false, "print result without body")

	flag.Parse()

	urlString, err := getUrl()
	if err != nil {
		return nil, err
	}

	var method string = "GET"
	if *isGet {
		method = "GET"
	} else if *isPost {
		method = "POST"
	} else if *isHead {
		method = "HEAD"
	} else if *customMethod != "" {
		method = *customMethod
	}

	return &UriParams{
		Url:       urlString,
		Method:    method,
		PrintBody: !*noBody,
	}, nil
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := NewRequest(&ConsolePrint{})
	err = r.Do(uriParams)
	if err != nil {
		fmt.Println(err)
		return
	}
}
