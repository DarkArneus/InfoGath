package main

import (
	"fmt"
	"sync"
	
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
)

func detectInput(domains []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, domain := range domains {
		c := colly.NewCollector()
		// When you find input
		c.OnHTML("input", func(e *colly.HTMLElement) {
			inputType := e.Attr("type")
			inputName := e.Attr("name")
			results <-fmt.Sprintf("%s\n  Input Type: %s\n  Input Name: %s\n", color.YellowString("Domain: "+domain), inputType, inputName)
		})
		// When you find forms
		c.OnHTML("form", func(e *colly.HTMLElement) {
			formAction := e.Attr("action")
			results <- fmt.Sprintf("%s\n  Form Action: %s\n", color.MagentaString("Domain: "+domain), formAction)
		})
		// Visits the domain
		err := c.Visit("http://" + domain)
		if err != nil {
			fmt.Printf("%s: %v\n", color.RedString("Error visiting domain "+domain), err)
		}
	}
}
