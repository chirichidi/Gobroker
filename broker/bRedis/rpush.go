package bRedis

import (
	"dancechanlibrary/broker"
)

type RPush struct {
	Key        string
	Value      string                // TODO 성능때문에 single parameter 사용
	ResultChan chan broker.CmdResult `json:"-"`
}

func (rp *RPush) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)

	result, err := r.client.RPush(rp.Key, rp.Value).Result()
	select {
	case rp.ResultChan <- broker.CmdResult{Result: result, Err: err}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (rp *RPush) Result() broker.CmdResult {
	return <-rp.ResultChan
}
