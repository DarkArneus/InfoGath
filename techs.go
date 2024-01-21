package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func init() {
	flag.Usage = func() {
		h := []string{
			"InfoGath. Gather all the information as fast as possible.",
			"",
			"Options:",
			"  -f, --file <txtfile>     Specify the URLs to fetch",
			"  -t, --threads <int>      Indicate the number of threads you want to use. Number of threads must be lower than number of domains!\n",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func parseTXT(domains string) []string {
	file, err := os.Open(domains)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
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

func statusCode(domains []string, index int, num_it int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := index; i < index+num_it; i++ {
		resp, err := http.Get("http://" + domains[i])
		if err != nil {
			results <- fmt.Sprintf("%s: Error - %s", domains[i], err.Error())
		} else {
			results <- fmt.Sprintf("%s: Status - %s", domains[i], resp.Status)
			resp.Body.Close()
		}
	}

}

func main() {
	var domainsFile string
	flag.StringVar(&domainsFile, "file", "", "Specify the file containing URLs to fetch")
	flag.StringVar(&domainsFile, "f", "", "Specify the file containing URLs to fetch (shorthand)")

	var threads int
	flag.IntVar(&threads, "threads", 1, "Indicate the number of threads you want to use")
	flag.IntVar(&threads, "t", 1, "Indicate the number of threads you want to use (shorthand)")

	flag.Parse()

	if domainsFile == "" {
		flag.Usage()
		return
	}
	domains := parseTXT(domainsFile)
	if len(domains) < threads {
		fmt.Println("Please use a lower number of threads ")
		flag.Usage()
		return
	}
	var iterate = len(domains) / threads // how many iterations x goroutine must be made

	var wg sync.WaitGroup
	results := make(chan string, len(domains))

	// Inicia workers
	for i := 0; i < iterate; i++ {
		wg.Add(1)
		var total_iterate = i + iterate
		go statusCode(domains, total_iterate, iterate, results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

}
