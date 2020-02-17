package bCurl

import (
	"net/http"
)

type bHttp struct {
	client *http.Client
}

func InitBHttp(client *http.Client) *bHttp {
	// 0. check
	//client.

	// 1. register
	_client := &bHttp{client: client}
	return _client
}

func (r *bHttp) Clear() {
	r.client = nil
}

func (r *bHttp) Ping() bool {
	return true
}
