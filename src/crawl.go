package main

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	//"github.com/gocolly/colly/v2/debug" for debugging
)

func visitAnchor(domains []string, maxDepth int, results chan<- string, wg *sync.WaitGroup){
	defer wg.Done()
	next_visit := make(map[string]bool)
	previousNode := make(map[string][]string)
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth), 
		//colly.Async(true),
		//colly.Debugger(&debug.LogDebugger{}),
	)
	stop := false

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		nextlink:= e.Request.AbsoluteURL(e.Attr("href"))
		if !next_visit[nextlink] && nextlink != "" && !contains(previousNode[nextlink], e.Request.URL.String()){
			// Mark as visited
			next_visit[nextlink] = true
			// Add the current node to the previous nodes
			previousNode[nextlink] = append(previousNode[nextlink], e.Request.URL.String())
			// Write in the results channel
			fmt.Printf("Visiting: %s that comes from: [%s]\n", color.YellowString(nextlink), color.CyanString(e.Request.URL.String()))
			results <- fmt.Sprintf("Visiting: %s that comes from: [%s]", color.YellowString(nextlink), color.CyanString(e.Request.URL.String()))
			e.Request.Visit(nextlink)
		}
	})	

	c.Limit(&colly.LimitRule{
		//Parallelism: 2,
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
	//c.Wait()
}

func crawlDetect(domains []string, results chan<- string, wg *sync.WaitGroup){
	
}

func contains(slice []string, value string) bool{
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}