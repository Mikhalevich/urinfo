package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type UriParams struct {
	Url       string
	Method    string
	PrintBody bool
}

func parseArguments() (*UriParams, error) {
	urlString := flag.String("url", "http://localhost:8080", "requesting url")
	isGet := flag.Bool("get", false, "get method")
	isPost := flag.Bool("post", false, "post method")
	isHead := flag.Bool("head", false, "head method")
	noBody := flag.Bool("nobody", false, "print result without body")

	flag.Parse()

	if len(*urlString) <= 0 {
		return nil, errors.New("Plese specify url")
	}

	uri, err := url.Parse(*urlString)
	if err != nil {
		return nil, err
	}

	if len(uri.Scheme) <= 0 {
		*urlString = "http://" + *urlString
	}

	var method string = "GET"
	if *isGet {
		method = "GET"
	} else if *isPost {
		method = "POST"
	} else if *isHead {
		method = "HEAD"
	}

	return &UriParams{
		Url:       *urlString,
		Method:    method,
		PrintBody: !*noBody,
	}, nil
}

func print(description string, status string, method string, headers *http.Header, body *io.ReadCloser) {
	fmt.Println(description)

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
	print("<<<<<<<<<<<<<<<<<<<<<<<< redirect", request.Response.Status, "", &request.Response.Header, nil)

	return nil
}

func doRequest(params *UriParams) {
	request, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	print(">>>>>>>>>>>>>>>>>>>>>>>> request", "", request.Method, &request.Header, nil)

	client := &http.Client{
		CheckRedirect: doRedirect,
	}

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

	print("<<<<<<<<<<<<<<<<<<<<<<<< result", response.Status, "", &response.Header, body)
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	doRequest(uriParams)
}
