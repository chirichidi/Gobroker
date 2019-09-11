package bRedis

import (
	"dancechanlibrary/broker"
)

type Get struct {
	Key        string
	ResultChan chan broker.CmdResult `json:"-"`
}

func (g *Get) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)
	result, err := r.client.Get(g.Key).Result()

	select {
	case g.ResultChan <- broker.CmdResult{Result: result, Err: err}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (g *Get) Result() broker.CmdResult {
	return <-g.ResultChan
}
