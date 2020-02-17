package bCurl

import "net/http"

type BHttpResponse struct {
	Header     http.Header
	Body       string
	Status     string
	StatusCode int
}
