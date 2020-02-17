package bCurl

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"testing"
)

func HttpMock() *http.Client {
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	return client
}

func Test_bHttp_생성_테스트(t *testing.T) {
	//given
	c := HttpMock()

	//when
	r := InitBHttp(c)
	b := r.Ping()

	//then
	assert.Equal(t, b, true)
}
