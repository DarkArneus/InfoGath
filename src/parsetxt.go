package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseTXT(domains string) []string {
	file, err := os.Open(domains)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	var array_domain []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		array_domain = append(array_domain, domain)
	}

	return array_domain
}

func parseCrawlTXT(domains string) []string {
    file, err := os.Open(domains)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return nil
    }
    defer file.Close()

    var array_domain []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        domain := scanner.Text()
        colonIndex := strings.Index(domain, ":")
        if colonIndex != -1 {
            domain = domain[:colonIndex]
            domain = strings.TrimSpace(domain)
            array_domain = append(array_domain, domain) 
        } else {
            fmt.Println(": Not found")
        }
    }
    return array_domain 
}
