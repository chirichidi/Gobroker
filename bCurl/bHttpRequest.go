package bCurl

import (
	"bytes"
	"dancechanlibrary/broker"
	"net/http"
)

type BHttpRequestMethod int

const (
	BHTTP_REQUEST_TYPE_GET  = BHttpRequestMethod(1)
	BHTTP_REQUEST_TYPE_POST = BHttpRequestMethod(2)
)

type bHttpRequest struct {
	request *http.Request
	option  BHttpOption
}

type BHttpOption struct {
	Header    map[string]string
	Url       string
	Method    BHttpRequestMethod
	PostField []byte
}

func MakeBHttpRequest(option BHttpOption) (*bHttpRequest, broker.ErrorCodeType) {

	// 0. check
	if option.Url == "" {
		return nil, broker.BROKER_ERROR_CODE_BCURL_HTTP_URL_NOT_SET
	}
	if option.Method == BHTTP_REQUEST_TYPE_POST {
		if option.PostField == nil {
			return nil, broker.BROKER_ERROR_CODE_BCURL_HTTP_POST_FIELD_NOT_SET
		}
	}

	// 1. make Request
	var httpRequest *http.Request
	var err error

	switch option.Method {
	case BHTTP_REQUEST_TYPE_GET:
		httpRequest, err = http.NewRequest("GET", option.Url, nil)
	case BHTTP_REQUEST_TYPE_POST:
		buffer := bytes.NewBuffer(option.PostField)
		httpRequest, err = http.NewRequest("POST", option.Url, buffer)
	default:
		return nil, broker.BROKER_ERROR_CODE_BCURL_HTTP_METHOD_NOT_SET
	}

	if err != nil {
		return nil, broker.BROKER_ERROR_CODE_MAKE_REQUEST_FAIL
	}

	if option.Header != nil {
		for key, value := range option.Header {
			httpRequest.Header.Add(key, value)
		}
	}

	request := &bHttpRequest{
		request: httpRequest,
		option:  option,
	}
	return request, broker.BROKER_ERROR_CODE_SUCCESS
}
