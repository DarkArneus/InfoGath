package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	flag.Usage = func() {
		h := []string{
			"InfoGath. Gather all the information as fast as possible.",
			"",
			"Options:",
			"  -f, --file <txtfile>     Specify the URLs to fetch",
			"  -t, --threads <int>      Indicate the number of threads you want to use\n",
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

func statusCode(url string) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	fileContent := fmt.Sprintf("%s: [%d]\n", url, resp.StatusCode)
	err = ioutil.WriteFile("output.txt", []byte(fileContent), 0644)
	if err != nil {
		log.Fatal("Error al escribir en el archivo:", err)
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
	for _, domain := range domains {
		statusCode(domain)
	}
}
