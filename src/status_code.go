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

        statusColor := color.GreenString // Por defecto, el color es verde
        switch resp.StatusCode {
        case 404:
            statusColor = color.RedString
        case 302:
            statusColor = color.BlueString
        case 403:
            statusColor = color.RedString
        }

        // Leer el título de la página
        title, err := getPageTitle(resp.Body)
        resp.Body.Close()

        if err != nil {
            results <- fmt.Sprintf("%s: Status - [%s], Error getting title", domain, statusColor(fmt.Sprintf("%s", resp.Status)))
        } else {
            results <- fmt.Sprintf("%s: Status - [%s] [%s]", domain, statusColor(fmt.Sprintf("%s", resp.Status)), color.MagentaString(title))
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
