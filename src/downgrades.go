package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/fatih/color"
)

func detectHttpDowngrades(domains []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Disable HTTP/2 by setting TLSNextProto to an empty map
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{
		Transport: transport,
	}

	// Loop through each domain and check for HTTP/1.1 support
	for _, domain := range domains {
		resp, err := client.Get("http://" + domain)
		if err != nil {
			results <- fmt.Sprintf("%s: %s", domain, color.RedString("[-] HTTP/1.1 Not Accepted - "+err.Error()))
			continue
		}

		// Close the response body explicitly
		if resp.Body != nil {
			resp.Body.Close()
		}

		// Check the protocol used
		if resp.Proto == "HTTP/1.1" {
			results <- fmt.Sprintf("%s: %s", domain, color.GreenString("[+] HTTP/1.1 Accepted"))
		} else {
			results <- fmt.Sprintf("%s: %s", domain, color.YellowString("[+] Downgraded successfully to HTTP/1.1"))
		}
	}
}
