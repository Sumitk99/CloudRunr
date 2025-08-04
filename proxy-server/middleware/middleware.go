package middleware

import (
	"github.com/Sumitk99/CloudRunr/proxy-server/constants"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//func ProxyHandler(c *gin.Context) {
//	hostname := c.Request.Host
//	parts := strings.Split(hostname, ".")
//	if len(parts) < 1 {
//		c.String(http.StatusBadRequest, "Invalid host")
//		return
//	}
//	subdomain := parts[0]
//
//	target, err := url.Parse(fmt.Sprintf("%s/%s/", constants.BASE_PATH, subdomain))
//	if err != nil {
//		c.String(http.StatusInternalServerError, "Invalid target URL")
//		return
//	}
//
//	proxy := httputil.NewSingleHostReverseProxy(target)
//
//	// Rewrite the request URL before proxying
//	originalDirector := proxy.Director
//	proxy.Director = func(req *http.Request) {
//		originalDirector(req)
//
//		// Fix path for root requests
//		if req.URL.Path == "/" {
//			req.URL.Path = "/index.html"
//		}
//	}
//
//	// Optional error handler
//	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
//		http.Error(rw, "Proxy error: "+err.Error(), http.StatusBadGateway)
//	}
//
//	proxy.ServeHTTP(c.Writer, c.Request)
//	c.Abort()
//}

func ProxyHandler(c *gin.Context) {
	hostname := c.Request.Host
	parts := strings.Split(hostname, ".")
	if len(parts) < 1 {
		c.String(http.StatusBadRequest, "Invalid host")
		return
	}
	subdomain := parts[0]

	// Build the base S3 folder URL
	targetURLStr := constants.BASE_PATH + "/" + subdomain
	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		c.String(http.StatusInternalServerError, "Invalid target URL")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Rewrite request path: "/" → "/index.html"
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Rewrite base path to include subdomain folder
		// Remove double slash if present
		if req.URL.Path == "/" {
			req.URL.Path = "/index.html"
		}

		// Ensure path doesn't prefix twice (safety)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
		req.URL.Path = "/" + req.URL.Path
	}

	// Optional: Proxy error handling
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		http.Error(rw, "Proxy error: "+err.Error(), http.StatusBadGateway)
	}

	// Serve the proxy
	proxy.ServeHTTP(c.Writer, c.Request)
	c.Abort()
}
