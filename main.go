package main

import (
	"errors"
	"flag"
	"fmt"
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

func printResult(params *UriParams, response *http.Response) {
	defer response.Body.Close()

	fmt.Printf("Status = %s\n", response.Status)
	fmt.Printf("Content-Lenght = %d\n", response.ContentLength)

	fmt.Println("*********** headers *************")
	for key, value := range response.Header {
		fmt.Printf("%s => %s\n", key, value)
	}

	if params.PrintBody {
		fmt.Println("*********** body ****************")
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))
	}
}

func doRequest(params *UriParams) {
	request, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	printResult(params, response)
}

func main() {
	uriParams, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	doRequest(uriParams)
}
