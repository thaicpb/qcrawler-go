package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Go Web Crawler v1.0")
	fmt.Println("===================")
	
	// Command line flags
	var (
		maxDepth     = flag.Int("depth", 2, "Maximum crawl depth")
		concurrency  = flag.Int("concurrency", 5, "Number of concurrent goroutines")
		userAgent    = flag.String("user-agent", "GoWebCrawler/1.0", "User-Agent header")
		timeout      = flag.Int("timeout", 10, "Timeout for each request (seconds)")
		configFile   = flag.String("config", "", "Path to config file (JSON)")
		selectors    = flag.String("selectors", "", "CSS selectors in format: name1=selector1,name2=selector2")
		interactive  = flag.Bool("interactive", false, "Interactive mode to input CSS selectors")
		outputFile   = flag.String("output", "output.json", "Output file name")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <url>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -selectors=\"title=h1,content=.article\" https://example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -config=config.json https://example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -interactive https://example.com\n", os.Args[0])
	}
	
	flag.Parse()
	
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	
	targetURL := flag.Arg(0)
	
	// Load configuration
	config := Config{
		MaxDepth:     *maxDepth,
		Concurrency:  *concurrency,
		UserAgent:    *userAgent,
		Timeout:      *timeout,
		CSSSelectors: make(map[string]string),
	}
	
	// Load from config file if provided
	if *configFile != "" {
		loadedConfig, err := loadConfigFromFile(*configFile)
		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}
		config = loadedConfig
	}
	
	// Override with command line selectors
	if *selectors != "" {
		parsedSelectors, err := parseSelectorsFlag(*selectors)
		if err != nil {
			log.Fatalf("Error parsing selectors: %v", err)
		}
		for k, v := range parsedSelectors {
			config.CSSSelectors[k] = v
		}
	}
	
	// Interactive mode
	if *interactive {
		interactiveSelectors := getInteractiveSelectors()
		for k, v := range interactiveSelectors {
			config.CSSSelectors[k] = v
		}
	}
	
	// Default selectors if none provided
	if len(config.CSSSelectors) == 0 {
		config.CSSSelectors = map[string]string{
			"headings":   "h1, h2, h3",
			"paragraphs": "p",
			"links":      "a",
		}
		fmt.Println("\nUsing default CSS selectors (headings, paragraphs, links)")
	}
	
	fmt.Printf("\nCSS Selectors configured:\n")
	for name, selector := range config.CSSSelectors {
		fmt.Printf("  %s: %s\n", name, selector)
	}
	
	crawler := NewCrawler(config)
	
	fmt.Printf("\nStarting crawl: %s\n", targetURL)
	
	results, err := crawler.Start(targetURL)
	if err != nil {
		log.Fatalf("Error while crawling: %v", err)
	}
	
	fmt.Printf("\nSuccessfully crawled %d pages\n", len(results))
	
	err = SaveResults(results, *outputFile)
	if err != nil {
		log.Printf("Error saving results: %v", err)
	} else {
		fmt.Printf("Results saved to %s\n", *outputFile)
	}
}

// ConfigFile represents the JSON config file structure
type ConfigFile struct {
	MaxDepth     int               `json:"max_depth"`
	Concurrency  int               `json:"concurrency"`
	UserAgent    string            `json:"user_agent"`
	Timeout      int               `json:"timeout"`
	CSSSelectors map[string]string `json:"css_selectors"`
}

func loadConfigFromFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	
	var configFile ConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return Config{}, err
	}
	
	return Config{
		MaxDepth:     configFile.MaxDepth,
		Concurrency:  configFile.Concurrency,
		UserAgent:    configFile.UserAgent,
		Timeout:      configFile.Timeout,
		CSSSelectors: configFile.CSSSelectors,
	}, nil
}

func parseSelectorsFlag(selectorsStr string) (map[string]string, error) {
	selectors := make(map[string]string)
	
	if selectorsStr == "" {
		return selectors, nil
	}
	
	pairs := strings.Split(selectorsStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid selector format: %s (expected name=selector)", pair)
		}
		
		name := strings.TrimSpace(kv[0])
		selector := strings.TrimSpace(kv[1])
		
		if name == "" || selector == "" {
			return nil, fmt.Errorf("empty name or selector in: %s", pair)
		}
		
		selectors[name] = selector
	}
	
	return selectors, nil
}

func getInteractiveSelectors() map[string]string {
	selectors := make(map[string]string)
	
	fmt.Println("\nInteractive CSS Selector Configuration")
	fmt.Println("=====================================")
	fmt.Println("Enter CSS selectors (press Enter without input to finish)")
	fmt.Println("Format: name=selector (e.g., title=h1)")
	fmt.Println()
	
	for {
		fmt.Print("CSS Selector: ")
		var input string
		fmt.Scanln(&input)
		
		if input == "" {
			break
		}
		
		kv := strings.SplitN(input, "=", 2)
		if len(kv) != 2 {
			fmt.Println("Invalid format. Please use: name=selector")
			continue
		}
		
		name := strings.TrimSpace(kv[0])
		selector := strings.TrimSpace(kv[1])
		
		if name == "" || selector == "" {
			fmt.Println("Name and selector cannot be empty")
			continue
		}
		
		selectors[name] = selector
		fmt.Printf("Added: %s = %s\n", name, selector)
	}
	
	return selectors
}