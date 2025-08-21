package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	PORT      = 8000
	BASE_PATH = "https://cloudrunr.s3.ap-south-1.amazonaws.com"
)

type SPAProxy struct {
	basePath   string
	httpClient *http.Client
}

func NewSPAProxy(basePath string) *SPAProxy {
	return &SPAProxy{
		basePath: basePath,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (sp *SPAProxy) extractSubdomain(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func (sp *SPAProxy) isStaticFile(path string) bool {
	ext := filepath.Ext(path)
	staticExts := []string{".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".woff", ".woff2", ".ttf", ".eot", ".map"}

	for _, staticExt := range staticExts {
		if ext == staticExt {
			return true
		}
	}
	return false
}

func (sp *SPAProxy) getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}

func (sp *SPAProxy) fetchFromS3(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return sp.httpClient.Do(req)
}

func (sp *SPAProxy) serveFile(c *gin.Context, s3URL string, originalPath string) {
	resp, err := sp.fetchFromS3(s3URL)
	if err != nil {
		log.Printf("Error fetching from S3: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch resource"})
		return
	}
	defer resp.Body.Close()

	// Copy headers from S3 response (except some)
	for key, values := range resp.Header {
		if key != "Content-Length" { // Let Gin handle this
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	// Set appropriate content type
	if c.GetHeader("Content-Type") == "" {
		c.Header("Content-Type", sp.getContentType(originalPath))
	}

	// Set status and copy body
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (sp *SPAProxy) handleRequest(c *gin.Context) {
	subdomain := sp.extractSubdomain(c.Request.Host)
	if subdomain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subdomain"})
		return
	}

	requestPath := c.Request.URL.Path
	log.Printf("Serving %s%s for subdomain: %s", c.Request.Host, requestPath, subdomain)

	var s3URL string

	if requestPath == "/" {
		// Root path - serve index.html
		s3URL = fmt.Sprintf("%s/%s/index.html", sp.basePath, subdomain)
		sp.serveFile(c, s3URL, "index.html")
	} else if sp.isStaticFile(requestPath) {
		// Static file - serve directly
		s3URL = fmt.Sprintf("%s/%s%s", sp.basePath, subdomain, requestPath)
		sp.serveFile(c, s3URL, requestPath)
	} else {
		// SPA route (like /orders/order_id) - serve index.html for client-side routing
		s3URL = fmt.Sprintf("%s/%s/index.html", sp.basePath, subdomain)

		// First try to fetch index.html
		resp, err := sp.fetchFromS3(s3URL)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		defer resp.Body.Close()

		// Serve index.html with correct content type
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusOK)
		io.Copy(c.Writer, resp.Body)
	}
}

func main() {
	r := gin.Default()

	// Initialize SPA proxy
	proxy := NewSPAProxy(BASE_PATH)

	// Middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s\" %d %s\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))
	r.Use(gin.Recovery())

	// Handle all routes
	r.NoRoute(proxy.handleRequest)

	fmt.Printf("🚀 SPA Proxy Server running on port %d\n", PORT)
	fmt.Printf("📁 Serving from: %s\n", BASE_PATH)
	fmt.Println("✅ Ready to serve SPAs with client-side routing!")

	err := r.Run(fmt.Sprintf("0.0.0.0:8000"))
	if err != nil {
		log.Fatal(err)
	}
}
