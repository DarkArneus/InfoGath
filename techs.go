package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

func init() {
	flag.Usage = func() {
		h := []string{
			"InfoGath. Gather all the information as fast as possible.",
			"",
			"Options:",
			"  -f, --file <txtfile>     Specify the URLs to fetch",
			"  -t, --threads <int>      Indicate the number of threads you want to use. Number of threads must be lower than number of domains!",
			"  -o, --output <file>		Indicate the name of the output file\n",
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

func main() {
	startTime := time.Now()
	var wg sync.WaitGroup

	var domainsFile string
	flag.StringVar(&domainsFile, "file", "", "Specify the file containing URLs to fetch")
	flag.StringVar(&domainsFile, "f", "", "Specify the file containing URLs to fetch (shorthand)")

	var threads int
	flag.IntVar(&threads, "threads", 1, "Indicate the number of threads you want to use")
	flag.IntVar(&threads, "t", 1, "Indicate the number of threads you want to use (shorthand)")

	var output string
	flag.StringVar(&output, "output", "", "Indicate the name of the output file")
	flag.StringVar(&output, "o", "", "Specify the file containing URLs to fetch (shorthand)")

	flag.Parse()

	if domainsFile == "" {
		flag.Usage()
		return
	}

	domains := parseTXT(domainsFile)
	results := make(chan string, len(domains))

	if len(domains) < threads {
		fmt.Println("Please use a lower number of threads ")
		flag.Usage()
		return
	}

	if threads == 1 {
		fmt.Println("Using default number of threads: len(domains) / 2")
		threads = len(domains) / 2
	}

	var domainsPerThread = len(domains) / threads // how many iterations x goroutine must be made

	// Inicia workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		start := i * domainsPerThread
		end := (i + 1) * domainsPerThread

		// For the last goroutine, include any remaining domains
		if i == threads-1 {
			end = len(domains)
		}

		go statusCode(domains[start:end], results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	if output == "" {
		output = "active_subdomains.txt"
	}

	// Create a file to write active subdomains
	outputFile, err := os.Create(output)
	defer outputFile.Close()
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}

	// Write active subdomains to the file
	for result := range results {
		//fmt.Println(result)                     // Print to console as before
		if !strings.Contains(result, "Error") { // Check if the status code indicates success (adjust this condition as needed)
			mu.Lock() // Acquire the lock before writing
			_, err := outputFile.WriteString(result + "\n")
			mu.Unlock() // Release the lock after writing
			if err != nil {
				fmt.Println("Error writing to output file:", err)
			}
		}
	}
	// Calculate and print the runtime
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Total runtime:", elapsedTime)

}
