package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UriParams struct {
	URL         string
	Method      string
	PrintBody   bool
	ForceHTTP11 bool
}

type Printer interface {
	Print(description string, delta time.Duration, total time.Duration, response *http.Response)
}

func getURL() (string, error) {
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
	http11 := flag.Bool("http11", false, "use HTTP/1.1 protocol")

	flag.Parse()

	urlString, err := getURL()
	if err != nil {
		return nil, err
	}

	method := "GET"
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
		URL:         urlString,
		Method:      method,
		PrintBody:   !*noBody,
		ForceHTTP11: *http11,
	}, nil
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := NewRequest(NewConsolePrint(uriParams.PrintBody))
	err = r.Do(uriParams.Method, uriParams.URL, uriParams.ForceHTTP11)
	if err != nil {
		fmt.Println(err)
		return
	}
}
