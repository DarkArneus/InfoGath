package main

import (
	"bufio"
	"fmt"
	"os"
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
