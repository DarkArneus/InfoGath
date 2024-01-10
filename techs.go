package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "bufio"
    //"strings"
    "strconv"
 )

func parseTXT(domains string)[]string {
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

func statusCode(url string){
    
    domainsFile := os.Args[1]
    domains := parseTXT(domainsFile)
    fileContent := ""
    for _, domain := range domains{
        resp, err := http.Get(domain)
        fileContent += domain + ": ["
        fileContent +=  strconv.Itoa(resp.StatusCode) 
        fileContent += "]\n"
        if err != nil {
            log.Fatal(err)
        }
    }
    err := ioutil.WriteFile("output.txt", []byte(fmt.Sprintf("%s", fileContent)), 0644)
    if err != nil {
        log.Fatal("Error al escribir en el archivo:", err)
    }
}

func main() {
    if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go domains.txt")
		return
	}

	domainsFile := os.Args[1]
	domains := parseTXT(domainsFile)

	for _, domain := range domains {
		statusCode(domain)
	}
}

