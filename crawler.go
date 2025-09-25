package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Config struct {
	MaxDepth     int
	Concurrency  int
	UserAgent    string
	Timeout      int
	CSSSelectors map[string]string // map[fieldName]cssSelector
}

type Crawler struct {
	config      Config
	visited     map[string]bool
	visitedLock sync.RWMutex
	client      *http.Client
	sem         chan struct{}
}

type Page struct {
	URL           string            `json:"url"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	Links         []string          `json:"links"`
	CrawledAt     time.Time         `json:"crawled_at"`
	SelectorData  map[string]string `json:"selector_data,omitempty"` // CSS selector extracted data
}

func NewCrawler(config Config) *Crawler {
	return &Crawler{
		config:  config,
		visited: make(map[string]bool),
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		sem: make(chan struct{}, config.Concurrency),
	}
}

func (c *Crawler) Start(startURL string) ([]Page, error) {
	_, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	
	var results []Page
	var resultsLock sync.Mutex
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		pages := c.crawl(startURL, 0)
		resultsLock.Lock()
		results = append(results, pages...)
		resultsLock.Unlock()
	}()
	
	wg.Wait()
	return results, nil
}

func (c *Crawler) crawl(targetURL string, depth int) []Page {
	if depth > c.config.MaxDepth {
		return nil
	}
	
	c.visitedLock.Lock()
	if c.visited[targetURL] {
		c.visitedLock.Unlock()
		return nil
	}
	c.visited[targetURL] = true
	c.visitedLock.Unlock()
	
	c.sem <- struct{}{}
	defer func() { <-c.sem }()
	
	page, err := c.fetchPage(targetURL)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", targetURL, err)
		return nil
	}
	
	var results []Page
	results = append(results, page)
	
	var wg sync.WaitGroup
	var resultsLock sync.Mutex
	
	for _, link := range page.Links {
		absoluteURL := c.toAbsoluteURL(targetURL, link)
		if absoluteURL == "" {
			continue
		}
		
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			pages := c.crawl(url, depth+1)
			resultsLock.Lock()
			results = append(results, pages...)
			resultsLock.Unlock()
		}(absoluteURL)
	}
	
	wg.Wait()
	return results
}

func (c *Crawler) fetchPage(targetURL string) (Page, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return Page{}, err
	}
	
	req.Header.Set("User-Agent", c.config.UserAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return Page{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, err
	}
	
	page := Page{
		URL:       targetURL,
		CrawledAt: time.Now(),
	}
	
	c.parseHTML(string(body), &page)
	c.extractSelectorData(string(body), &page)
	
	return page, nil
}

func (c *Crawler) parseHTML(htmlContent string, page *Page) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return
	}
	
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					page.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						page.Links = append(page.Links, attr.Val)
					}
				}
			case "p", "h1", "h2", "h3", "h4", "h5", "h6":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					text := strings.TrimSpace(n.FirstChild.Data)
					if text != "" {
						page.Content += text + "\n"
					}
				}
			}
		}
		
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	
	extract(doc)
}

func (c *Crawler) toAbsoluteURL(baseURL, link string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	parsedLink, err := url.Parse(link)
	if err != nil {
		return ""
	}
	
	if parsedLink.IsAbs() {
		if parsedLink.Host != base.Host {
			return ""
		}
		return parsedLink.String()
	}
	
	return base.ResolveReference(parsedLink).String()
}

func (c *Crawler) extractSelectorData(htmlContent string, page *Page) {
	if len(c.config.CSSSelectors) == 0 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return
	}

	page.SelectorData = make(map[string]string)

	for fieldName, selector := range c.config.CSSSelectors {
		var extractedTexts []string
		
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				extractedTexts = append(extractedTexts, text)
			}
		})
		
		if len(extractedTexts) > 0 {
			page.SelectorData[fieldName] = strings.Join(extractedTexts, " | ")
		}
	}
}