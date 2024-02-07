package main

import (
	"fmt"
	"net/http"
	"sync"
	
	"github.com/fatih/color"
)

func statusCode(domains []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, domain := range domains {
		resp, err := http.Get("http://" + domain)
		if err != nil {
			results <- fmt.Sprintf("%s: %s", color.RedString(domain), color.RedString("Error - "+err.Error()))
		} else {
			statusColor := color.GreenString // Por defecto, el color es verde
			switch resp.StatusCode {
			case 404:
				statusColor = color.RedString
			case 302:
				statusColor = color.BlueString
			case 403:
				statusColor = color.RedString
			}
			
			results <- fmt.Sprintf("%s: Status - [%s]", color.CyanString(domain), statusColor(fmt.Sprintf("%s", resp.Status)))
			resp.Body.Close()
		}
	}
}
