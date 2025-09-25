# QCrawler - Advanced Web Crawler

🕷️ A powerful, flexible web crawler written in Go with advanced CSS selector support and multiple configuration options.

## ✨ Features

- **Advanced CSS Selectors**: Extract data using powerful CSS selector syntax
- **Multiple Configuration Methods**: Command line, config files, or interactive mode
- **Multi-threaded Crawling**: Configurable concurrency with rate limiting
- **Flexible Output**: JSON format with custom field extraction
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Domain-safe**: Only crawls within the same domain
- **Professional CLI**: Easy-to-use command line interface

## 🚀 Installation

### Quick Install (Recommended)

```bash
# Clone the repository
git clone https://github.com/thaicpb/qcrawler-go.git
cd qcrawler-go

# Build and install
make install
```

### Manual Installation

```bash
# Clone and build
git clone https://github.com/thaicpb/qcrawler-go.git
cd qcrawler-go
make build

# Copy to your PATH
sudo cp build/qcrawler /usr/local/bin/
```

### Development Setup

```bash
git clone https://github.com/thaicpb/qcrawler-go.git
cd qcrawler-go
make deps
make build
```

## 📖 Usage

### Basic Usage

```bash
# Simple crawl
qcrawler -u https://example.com

# With custom output file
qcrawler -u https://example.com -o results.json

# Verbose mode
qcrawler -u https://example.com -v
```

### Advanced Usage

```bash
# Custom CSS selectors
qcrawler -u https://example.com -s "title=h1,content=.article,author=.byline"

# Using config file
qcrawler -u https://example.com -f config.json

# Interactive mode
qcrawler -u https://example.com -i

# Custom depth and concurrency
qcrawler -u https://example.com -d 3 -c 10
```

## 🛠️ Configuration

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-u` | URL to crawl (required) | - |
| `-f` | Config file (JSON) | - |
| `-d` | Max crawl depth | 2 |
| `-c` | Concurrent workers | 5 |
| `-o` | Output file | output.json |
| `-t` | Timeout (seconds) | 10 |
| `-s` | CSS selectors | - |
| `-i` | Interactive mode | false |
| `-v` | Verbose output | false |
| `-h` | Show help | - |

### Config File Format

Create a `config.json` file:

```json
{
  "max_depth": 3,
  "concurrency": 8,
  "user_agent": "QCrawler/1.0",
  "timeout": 15,
  "css_selectors": {
    "main_title": "h1",
    "article_content": "article p, .content p",
    "author": ".author, .byline, [data-author]",
    "publish_date": ".date, time, .publish-date",
    "meta_description": "meta[name='description']",
    "navigation": "nav a, .navbar a",
    "tags": ".tag, .tags a, .categories a"
  }
}
```

### CSS Selector Examples

```bash
# Basic elements
qcrawler -u https://example.com -s "title=h1,paragraphs=p"

# Classes and IDs
qcrawler -u https://example.com -s "content=.article,sidebar=#sidebar"

# Attributes
qcrawler -u https://example.com -s "meta_desc=meta[name='description']"

# Complex selectors
qcrawler -u https://example.com -s "nav_links=nav a,author=.author .name"

# Multiple elements (combined with |)
qcrawler -u https://example.com -s "headings=h1,h2,h3,content=article p"
```

## 📁 Output Format

```json
[
  {
    "url": "https://example.com",
    "title": "Page Title",
    "content": "Extracted text content...",
    "links": ["https://example.com/page1"],
    "crawled_at": "2024-01-01T10:00:00Z",
    "selector_data": {
      "main_title": "Welcome to Example",
      "author": "John Doe",
      "meta_description": "This is an example website",
      "article_content": "Main article text here..."
    }
  }
]
```

## 🛠️ Development

### Build Commands

```bash
# Show all available commands
make help

# Development
make build          # Build binary
make dev            # Build and run example
make test           # Run tests
make fmt            # Format code

# Installation
make install        # Install to /usr/local/bin
make uninstall      # Remove from system

# Cross-platform
make build-all      # Build for all platforms
make release        # Create release archives
```

### Project Structure

```
qcrawler-go/
├── cmd/qcrawler/           # CLI application
│   └── main.go
├── internal/crawler/       # Core crawler logic
│   ├── crawler.go
│   └── storage.go
├── build/                  # Build output
├── config.json            # Example config
├── Makefile              # Build automation
├── go.mod                # Go module
└── README.md             # Documentation
```

## 🔧 Advanced Examples

### News Website Scraping

```json
{
  "max_depth": 2,
  "concurrency": 5,
  "css_selectors": {
    "headline": "h1, .headline",
    "article_body": "article p, .article-content p",
    "author": ".author, .byline",
    "publish_date": "time, .date, .publish-date",
    "category": ".category, .section",
    "tags": ".tags a, .tag",
    "related_links": ".related a"
  }
}
```

### E-commerce Site

```json
{
  "css_selectors": {
    "product_title": "h1.product-title, .product-name",
    "price": ".price, .cost, .product-price",
    "description": ".product-description, .description",
    "specs": ".specifications p, .specs li",
    "reviews": ".review-text, .review-content",
    "rating": ".rating, .stars"
  }
}
```

### Blog Scraping

```bash
qcrawler -u https://blog.example.com -s "title=.post-title,content=.post-content,author=.author-name,date=.post-date,tags=.tag-list a" -d 3 -v
```

## 🚀 Performance Tips

1. **Adjust Concurrency**: Use `-c` flag based on target website capacity
2. **Limit Depth**: Use `-d` to prevent excessive crawling
3. **Use Specific Selectors**: More specific CSS selectors are faster
4. **Enable Verbose Mode**: Use `-v` to monitor performance

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Submit a pull request

## 📜 License

This project is open source. See LICENSE file for details.

## 🆘 Support

- Run `qcrawler -h` for help
- Check the examples in this README
- Open an issue for bugs or feature requests

---

**Happy Crawling! 🕷️**