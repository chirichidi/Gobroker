package bCurl

import (
	"dancechanlibrary/broker"
	"io/ioutil"
)

type DoHttp struct {
	*bHttpRequest
	resultChan chan *BHttpResponse
	err        error
}

func (do *DoHttp) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bHttp)
	do.resultChan = make(chan *BHttpResponse, 1)
	defer func() {
		close(do.resultChan)
	}()

	res, err := r.client.Do(do.request)
	if err != nil {
		do.err = err
		return broker.BROKER_ERROR_CODE_BCURL_HTTP_RESPONSE_FAIL
	}
	defer func() {
		_ = res.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		do.err = err
		return broker.BROKER_ERROR_CODE_BCURL_HTTP_READ_RESPONSE_BODY_FAIL
	}

	result := &BHttpResponse{
		Header:     res.Header,
		Body:       string(respBody),
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}

	select {
	case do.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (do *DoHttp) Result() (*BHttpResponse, error) {
	for result := range do.resultChan { // 채널 닫힐 때까지 blocking
		return result, nil
	}
	return nil, do.err
}
