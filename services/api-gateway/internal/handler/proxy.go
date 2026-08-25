package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy returns an HTTP reverse proxy for the given target URL.
func NewReverseProxy(target string) http.Handler {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}
