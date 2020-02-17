package bCurl

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func bHttpDoMock(option BHttpOption) *DoHttp {
	request, _ := MakeBHttpRequest(option)

	return &DoHttp{
		bHttpRequest: request,
	}
}

func Test_Do_Get_메소드_동작_확인(t *testing.T) {
	//given
	do := bHttpDoMock(BHttpOption{
		Header:    nil,
		Url:       "http://www.google.com",
		Method:    BHTTP_REQUEST_TYPE_GET,
		PostField: nil,
	})
	w := InitBHttp(&http.Client{
		Timeout: time.Second,
	})

	//when
	do.Cmd(w)

	//then
	result, err := do.Result()
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
}

func Test_Do_Post_메소드_동작_확인(t *testing.T) {
	//given
	do := bHttpDoMock(BHttpOption{
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Url:       "https://test-dancechat.com2us.net:9090/gateway",
		Method:    BHTTP_REQUEST_TYPE_POST,
		PostField: make([]byte, 0),
	})
	w := InitBHttp(&http.Client{
		Timeout: time.Second,
	})

	//when
	do.Cmd(w)

	//then
	result, err := do.Result()
	assert.Equal(t, err, nil)
	assert.Equal(t, result.StatusCode, 200)
}

func Test_Do_유효하지_않은_url_동작_확인(t *testing.T) {
	//given
	do := bHttpDoMock(BHttpOption{
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Url:       "http://NOTEXISTPAGE.com",
		Method:    BHTTP_REQUEST_TYPE_POST,
		PostField: make([]byte, 0),
	})
	w := InitBHttp(&http.Client{
		Timeout: time.Second,
	})

	//when
	do.Cmd(w)

	//then
	result, err := do.Result()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
