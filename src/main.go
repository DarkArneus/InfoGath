package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/common-nighthawk/go-figure"
)

var mu sync.Mutex

func init() {
	myFigure := figure.NewFigure("InfoGath", "", true)
	myFigure.Print()
	flag.Usage = func() {
		h := []string{
			"\tGather all the information as fast as possible.",
			"",
			"Options:",
			"  -f, --file <input file>          Specify the URLs to fetch and return the status code",
			"  -t, --threads <int>              Indicate the number of threads you want to use. Number of threads must be lower than number of domains!",
			"  -o, --output <output file>	   Indicate the name of the output file",
			"  -d, --detect <input file>        Visit all anchors given in the input file",
			"  -c, --crawl	<input file>	   Indicate to detect inputs as forms or labels",
			"       --depth <int>		   Indicate depth for the crawl",
			"  -dw, --downgrade		   Check if we can downgrade to HTTP1.1\n",
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
	flag.StringVar(&crawl, "crawl", "", "Specificate a depth and crawl the subdomains given")
	flag.StringVar(&crawl, "c", "", "Specificate a depth and crawl the subdomains given (shorthand)")

	var depth int
	flag.IntVar(&depth, "depth", 0, "Indicate the depth for the crawler")

	var detect string
	flag.StringVar(&detect, "detect", "", "Indicate whether to detect inputs as forms or input labels")
	flag.StringVar(&detect, "d", "", "Indicate whether to detect inputs as forms or input labels (shorthand)")

	var downgrade string
	flag.StringVar(&downgrade, "downgrade", "", "Check if we can downgrade to HTTP/1.1")
	flag.StringVar(&downgrade, "dw", "", "Check if we can downgrade to HTTP/1.1 (shorthand)")

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
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
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

	if detect != "" {
		outputDetectFile, err := os.Create("detect.txt")
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputDetectFile.Close()
		detect_domains := parseDetectTXT(detect)
		detect_results := make(chan string, len(detect_domains))
		var detectsPerThread = len(detect_domains) / threads
		for i := 0; i < threads; i++ {
			wg.Add(1)
			start := i * detectsPerThread
			end := (i + 1) * detectsPerThread

			// For the last goroutine, include any remaining domains
			if i == threads-1 {
				end = len(detect_domains)
			}

			go detectInput(detect_domains[start:end], detect_results, &wg)
		}
		// Anonymous goroutine designed to wait for closing the channel
		go func() {
			wg.Wait()
			close(detect_results)
		}()

		for detect_result := range detect_results {
			mu.Lock() // Acquire the lock before writing
			_, errc := outputDetectFile.WriteString(detect_result + "\n")
			mu.Unlock() // Release the lock after writing
			if errc != nil {
				fmt.Println("Error writing to output file:", err)

			}
		}
	}

	// crawler
	if crawl != "" {
		outputCrawlFile, err := os.Create("crawl.txt")
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputCrawlFile.Close()
		crawl_domains := parseDetectTXT(crawl)
		crawl_results := make(chan string, len(crawl_domains))
		var crawlPerThread = len(crawl_domains) / threads
		for i := 0; i < threads; i++ {
			wg.Add(1)
			start := i * crawlPerThread
			end := (i + 1) * crawlPerThread

			// For the last goroutine, include any remaining domains
			if i == threads-1 {
				end = len(crawl_domains)
			}

			go visitAnchor(crawl_domains[start:end], depth, crawl_results, &wg)
		}

		go func() {
			wg.Wait()
			close(crawl_results)
		}()

		for crawl_result := range crawl_results {
			mu.Lock() // Acquire the lock before writing
			_, errc := outputCrawlFile.WriteString(crawl_result + "\n")
			mu.Unlock() // Release the lock after writing
			if errc != nil {
				fmt.Println("Error writing to output file:", err)

			}
		}
	}

	if downgrade != "" {
		outputDowngrades, err := os.Create("downgrades.txt")
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputDowngrades.Close()
	
		// Parsear los dominios y dividirlos por hilos
		downgradeDomains := parseTXT(downgrade)
		downgradeResults := make(chan string, len(downgradeDomains))
		downgradePerThread := len(downgradeDomains) / threads
	
		for i := 0; i < threads; i++ {
			wg.Add(1)
			start := i * downgradePerThread
			end := (i + 1) * downgradePerThread
	
			// Incluir los dominios restantes en el último hilo
			if i == threads-1 {
				end = len(downgradeDomains)
			}
	
			go detectHttpDowngrades(downgradeDomains[start:end], downgradeResults, &wg)
		}
	
		go func() {
			wg.Wait()
			close(downgradeResults)
		}()
	
		// Escribir los resultados en el archivo de salida
		for downgradeResult := range downgradeResults {
			mu.Lock() // Adquirir el lock antes de escribir
			_, errc := outputDowngrades.WriteString(downgradeResult + "\n")
			mu.Unlock() // Liberar el lock después de escribir
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
