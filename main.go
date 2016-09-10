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
	Url    string
	Method string
}

func parseArguments() (*UriParams, error) {
	urlString := flag.String("url", "http://localhost:8080", "requesting url")
	isGet := flag.Bool("get", false, "get method")
	isPost := flag.Bool("post", false, "post method")
	isHead := flag.Bool("head", false, "head method")

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

	var method string
	if *isGet {
		method = "GET"
	} else if *isPost {
		method = "POST"
	} else if *isHead {
		method = "HEAD"
	}

	if len(method) <= 0 {
		return nil, errors.New("Invalid http method")
	}

	return &UriParams{
		Url:    *urlString,
		Method: method,
	}, nil
}

func printResult(response *http.Response, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	fmt.Printf("Status = %s\n", response.Status)
	fmt.Printf("Content-Lenght = %d\n", response.ContentLength)

	fmt.Println("*********** headers *************")
	for key, value := range response.Header {
		fmt.Printf("%s => %s\n", key, value)
	}

	if response.Request.Method == "POST" {
		fmt.Println("*********** body ****************")
		fmt.Println(ioutil.ReadAll(response.Body))
	} else {
		io.Copy(ioutil.Discard, response.Body)
	}
}

func Get(url string) {
	response, err := http.Get(url)
	printResult(response, err)
}

func doRequest(params *UriParams) {
	request, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	printResult(response, err)
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	doRequest(uriParams)
}
