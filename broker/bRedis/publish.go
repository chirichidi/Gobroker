package bRedis

import (
	"dancechanlibrary/broker"
)

type Publish struct {
	Channel    string
	Message    string
	ResultChan chan broker.CmdResult `json:"-"`
}

func (p *Publish) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)

	result, err := r.client.Publish(p.Channel, p.Message).Result()
	select {
	case p.ResultChan <- broker.CmdResult{Result: result, Err: err}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (p *Publish) Result() broker.CmdResult {
	return <-p.ResultChan
}
