package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewReverseProxy returns an HTTP reverse proxy configured to strip the public
// API version prefix before forwarding the request to the downstream service.
func NewReverseProxy(target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid proxy target", http.StatusInternalServerError)
		})
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Director = func(r *http.Request) {
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
		r.Host = u.Host
		r.URL.Path = rewritePath(r.URL.Path)
		r.URL.RawPath = rewritePath(r.URL.EscapedPath())
	}

	return proxy
}

func rewritePath(path string) string {
	trimmed := strings.TrimPrefix(path, "/api/v1")
	if trimmed == "" {
		return "/"
	}
	return trimmed
}
