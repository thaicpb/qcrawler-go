package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Go Web Crawler v1.0")
	fmt.Println("===================")
	
	if len(os.Args) < 2 {
		log.Fatal("Please provide a URL to crawl. Usage: go run main.go <url>")
	}
	
	targetURL := os.Args[1]
	
	crawler := NewCrawler(Config{
		MaxDepth:    2,
		Concurrency: 5,
		UserAgent:   "GoWebCrawler/1.0",
		Timeout:     10,
	})
	
	fmt.Printf("\nStarting crawl: %s\n", targetURL)
	
	results, err := crawler.Start(targetURL)
	if err != nil {
		log.Fatalf("Error while crawling: %v", err)
	}
	
	fmt.Printf("\nSuccessfully crawled %d pages\n", len(results))
	
	err = SaveResults(results, "output.json")
	if err != nil {
		log.Printf("Error saving results: %v", err)
	} else {
		fmt.Println("Results saved to output.json")
	}
}