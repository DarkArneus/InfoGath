package main

import (
	"fmt"
	"net/http"
	"sync"
)

func statusCode(domains []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, domain := range domains {
		resp, err := http.Get("http://" + domain)
		if err != nil {
			results <- fmt.Sprintf("%s: Error - %s", domain, err.Error())
		} else {
			results <- fmt.Sprintf("%s: Status - %s", domain, resp.Status)
			resp.Body.Close()
		}
	}
}
