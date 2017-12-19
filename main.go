package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type UriParams struct {
	Url       string
	Method    string
	PrintBody bool
}

var (
	StartTime    time.Time
	PreviousTime time.Time
)

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

func print(description string, delta time.Duration, total time.Duration, status string, method string, headers *http.Header, body *io.ReadCloser) {
	fmt.Printf("%s delta = %s total = %s\n", description, delta, total)

	if len(status) > 0 {
		fmt.Printf("Status = %s\n", status)
	}

	if len(method) > 0 {
		fmt.Printf("Method = %s\n", method)
	}

	fmt.Println("HEADERS:")
	for key, value := range *headers {
		fmt.Printf("%s => %s\n", key, value)
	}

	if body != nil {
		fmt.Println("BODY:")
		body, err := ioutil.ReadAll(*body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))
	}
}

func doRedirect(request *http.Request, via []*http.Request) error {
	currentTime := time.Now()
	print("<<<<<<<<<<<<<<<<<<<<<<<< redirect", currentTime.Sub(PreviousTime), currentTime.Sub(StartTime), request.Response.Status, "", &request.Response.Header, nil)
	PreviousTime = currentTime

	return nil
}

func doRequest(params *UriParams) {
	request, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{
		CheckRedirect: doRedirect,
	}

	StartTime = time.Now()
	PreviousTime = StartTime

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	var body *io.ReadCloser = nil
	if params.PrintBody {
		body = &response.Body
	}

	currentTime := time.Now()
	print("<<<<<<<<<<<<<<<<<<<<<<<< result", currentTime.Sub(PreviousTime), currentTime.Sub(StartTime), response.Status, "", &response.Header, body)
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	doRequest(uriParams)
}
