package bRedis

import (
	"dancechanlibrary/broker"
)

type Subscribe struct {
	Channel    string                //TODO 성능 이슈로 한개만
	ResultChan chan broker.CmdResult `json:"-"`
}

func (s *Subscribe) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)
	result := r.client.Subscribe(s.Channel)

	select {
	case s.ResultChan <- broker.CmdResult{Result: result}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (s *Subscribe) Result() broker.CmdResult {
	return <-s.ResultChan
}
