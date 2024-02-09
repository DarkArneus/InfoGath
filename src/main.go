package main

import (
	"flag"
	"fmt"
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
			"  -f, --file <file>     Specify the URLs to fetch",
			"  -t, --threads <int>      Indicate the number of threads you want to use. Number of threads must be lower than number of domains!",
			"  -o, --output <file>	   Indicate the name of the output file",
			"  -c, --crawl	<file>	   Indicate to detect inputs as forms or labels\n",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
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
	flag.StringVar(&output, "output", "active_subdomains", "Indicate the name of the output file")
	flag.StringVar(&output, "o", "active_subdomains", "Specify the file containing URLs to fetch (shorthand)")

	var crawl string
	flag.StringVar(&crawl, "crawl", "", "Indicate whether to detect inputs as forms or input labels")
	flag.StringVar(&crawl, "c", "", "Indicate whether to detect inputs as forms or input labels (shorthand)")

	flag.Parse()

	if domainsFile == "" {
		flag.Usage()
		return
	}

	domains := parseTXT(domainsFile)
	results := make(chan string, len(domains))

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

	// Create a file to write active subdomains
	outputFile, err := os.Create(output)
	defer outputFile.Close()
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}

	// Write active subdomains to the file
	for result := range results {
		if !strings.Contains(result, "Error") { // Check if the status code indicates success (adjust this condition as needed)
			mu.Lock() // Acquire the lock before writing
			_, err := outputFile.WriteString(result + "\n")
			mu.Unlock() // Release the lock after writing
			if err != nil {
				fmt.Println("Error writing to output file:", err)
			}
		}
	}

	if crawl != "" {
		outputCrawlFile, err := os.Create("crawl.txt")
		defer outputCrawlFile.Close()
		crawl_domains := parseCrawlTXT(crawl)
		crawl_results := make(chan string, len(crawl_domains))
		var crawlsPerThread = len(crawl_domains) / threads
		for i := 0; i < threads; i++ {
			wg.Add(1)
			start := i * crawlsPerThread
			end := (i + 1) * crawlsPerThread
	
			// For the last goroutine, include any remaining domains
			if i == threads-1 {
				end = len(crawl_domains)
			}
	
			go detectInput(crawl_domains[start:end], crawl_results, &wg)
		}
		
		go func() {
			wg.Wait()
			close(crawl_results)
			}()
			
			if crawl != ""{
				for crawl_result := range crawl_results {
					mu.Lock() // Acquire the lock before writing
					_, errc := outputCrawlFile.WriteString(crawl_result + "\n")
					mu.Unlock() // Release the lock after writing
					if errc != nil {
						fmt.Println("Error writing to output file:", err)
					}
				}
			}
	
	// Calculate and print the runtime
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Total runtime:", elapsedTime)

	}
}