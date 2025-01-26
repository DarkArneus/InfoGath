package main

import (
	"fmt"
    "net/http"
    "sync"
    "io"      // Para trabajar con la interfaz de E/S
    "strings" // Para manipular cadenas
    "errors"  // Para manejar errores
	
	"golang.org/x/net/html"
	"github.com/fatih/color"
)

func statusCode(domains []string, results chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    for _, domain := range domains {
        resp, err := http.Get("http://" + domain)
        if err != nil {
            results <- fmt.Sprintf("%s: %s", domain, color.RedString("Error - "+err.Error()))
            continue
        }

        statusColor := color.GreenString // Default green
        switch resp.StatusCode {
        case 404:
            statusColor = color.RedString
        case 302:
            statusColor = color.BlueString
        case 403:
            statusColor = color.RedString
        }

        // Read page title
        title, err := getPageTitle(resp.Body)
        cache := detectCache(resp)
        resp.Body.Close()

        var cacheString string
        if cache {
            cacheString = color.GreenString("Yes") // Cache found
        } else {
            cacheString = color.RedString("No") // No cache
        }

        if err != nil {
            results <- fmt.Sprintf("%s: Status - [%s], Error getting title. Cache - [%s]", domain, statusColor(fmt.Sprintf("%s", resp.Status)), cacheString)
        } else {
            results <- fmt.Sprintf("%s: Status - [%s] [%s] Cache - [%s]", domain, statusColor(fmt.Sprintf("%s", resp.Status)), color.MagentaString(title), cacheString)
        }
    }
}

// Obtain title
func getPageTitle(body io.Reader) (string, error) {
    doc, err := html.Parse(body)
    if err != nil {
        return "", err
    }
    titleNode := findTitleNode(doc)
    if titleNode == nil {
        return "", errors.New("Title not found")
    }
    return strings.TrimSpace(titleNode.FirstChild.Data), nil
}

// Finds title in DOM
func findTitleNode(n *html.Node) *html.Node {
    if n.Type == html.ElementNode && n.Data == "title" {
        return n
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if result := findTitleNode(c); result != nil {
            return result
        }
    }
    return nil
}

func detectCache(resp *http.Response) bool{
    if resp == nil{
        return false
    }

    cacheControl := resp.Header.Get("Cache-Control")
    expires := resp.Header.Get("Expires")
    eTag := resp.Header.Get("ETag")

    if cacheControl != "" || expires != "" || eTag != "" {
        return true
    }

    return false

}