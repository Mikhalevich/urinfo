package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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

func main() {
	Get("http://tut.by")
}
