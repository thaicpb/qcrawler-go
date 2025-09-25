# Go Web Crawler

A simple crawler tool written in Go for scraping static web pages.

## Features

- Multi-threaded crawling with concurrency limits
- Configurable crawl depth limit
- Extract title, content and links from HTML
- Save results in JSON format
- Automatically convert relative URLs to absolute
- Only crawl within the same domain

## Installation

```bash
go get golang.org/x/net/html
```

## Usage

```bash
go run . https://example.com
```

Results will be saved to `output.json`

## Configuration

You can adjust parameters in the `config.json` file:
- `max_depth`: Maximum crawl depth (default: 2)
- `concurrency`: Number of concurrent goroutines (default: 5)
- `user_agent`: User-Agent header
- `timeout`: Timeout for each request (seconds)

## Output Data Structure

```json
[
  {
    "url": "https://example.com",
    "title": "Example Domain",
    "content": "Extracted content...",
    "links": ["https://example.com/page1", "..."],
    "crawled_at": "2024-01-01T10:00:00Z"
  }
]
```