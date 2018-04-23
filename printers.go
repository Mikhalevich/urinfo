package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type EmptyPrint struct {
	// pass
}

func (cp *EmptyPrint) Print(description string, delta time.Duration, total time.Duration, status string, headers *http.Header, body *io.ReadCloser) {
	// pass
}

type ConsolePrint struct {
	// pass
}

func (cp *ConsolePrint) Print(description string, delta time.Duration, total time.Duration, status string, headers *http.Header, body *io.ReadCloser) {
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<< %s delta = %s total = %s\n", description, delta, total)

	if status != "" {
		fmt.Printf("Status = %s\n", status)
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
