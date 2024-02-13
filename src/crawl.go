package main

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	//"github.com/gocolly/colly/v2/debug"
)

func visitAnchor(domains []string, maxDepth int, results chan<- string, wg *sync.WaitGroup){
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth), 
		colly.Async(true),
		//colly.Debugger(&debug.LogDebugger{}),
	)
	stop := false

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		nextlink:= e.Request.AbsoluteURL(e.Attr("href"))
		results <-fmt.Sprintf("Visiting: %s that comes from: %s", color.YellowString(nextlink), e.Request.URL.String())
		e.Request.Visit(nextlink)
	})	

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	for _, domain := range domains {
		if stop != true{
			err := c.Visit("http://" + domain)
			if err != nil {
				fmt.Printf("%s: %v\n", color.RedString("Error visiting domain "+domain), err)
			}
		}

	}
	c.Wait()
}

func crawlDetect(domains []string, results chan<- string, wg *sync.WaitGroup){

}