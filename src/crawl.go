package main

import (
	"fmt"
	"sync"
	
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
)

func visitAnchor(domains []string, maxDepth int, results chan<- string, wg *sync.WaitGroup){
	for _, domain := range domains {
		c := colly.NewCollector(
			colly.MaxDepth(maxDepth), 
		)		
		
		c.OnHTML("a", func(e *colly.HTMLElement) {
			nextlink:= e.Request.AbsoluteURL(e.Attr("href"))
			c.Visit(nextlink)
			results <-fmt.Sprintf("Visiting: %s that comes from: %s", color.YellowString(nextlink), domain)
		})
		err := c.Visit("http://" + domain)
		if err != nil {
			fmt.Printf("%s: %v\n", color.RedString("Error visiting domain "+domain), err)
		}
	}
}

func crawlDetect(domains []string, results chan<- string, wg *sync.WaitGroup){

}