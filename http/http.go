package main

import (
	"fmt"
	"net/http"
)

func main() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for i, url := range urls {
				fmt.Println(i)

				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
			fmt.Println("chan done")
		}()
		fmt.Println("go call done")
		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://google.com", "https://badhost"}
	for responses := range checkStatus(done, urls...) {
		fmt.Printf("Respose %s\n", responses.Status)
	}
}
