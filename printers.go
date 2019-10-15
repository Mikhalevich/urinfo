package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type EmptyPrint struct {
	// pass
}

func (cp *EmptyPrint) Print(description string, delta time.Duration, total time.Duration, response *http.Response) {
	// pass
}

type ConsolePrint struct {
	PrintBody bool
}

func NewConsolePrint(pb bool) *ConsolePrint {
	return &ConsolePrint{
		PrintBody: pb,
	}
}

func (cp *ConsolePrint) Print(description string, delta time.Duration, total time.Duration, response *http.Response) {
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<< %s delta = %s total = %s\n", description, delta, total)

	fmt.Printf("Status: %s %s\n", response.Proto, response.Status)

	fmt.Println("HEADERS:")
	for key, value := range response.Header {
		fmt.Printf("%s => %s\n", key, value)
	}

	if !cp.PrintBody {
		return
	}

	if response.Body != nil {
		fmt.Println("BODY:")
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))
	}
}
