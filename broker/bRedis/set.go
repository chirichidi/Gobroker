package bRedis

import (
	"dancechanlibrary/broker"
	"time"
)

type Set struct {
	Key        string
	Value      string
	Time       time.Duration
	ResultChan chan broker.CmdResult `json:"-"`
}

func (s *Set) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)
	result, err := r.client.Set(s.Key, s.Value, s.Time).Result()

	select {
	case s.ResultChan <- broker.CmdResult{Result: result, Err: err}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (s *Set) Result() broker.CmdResult {
	return <-s.ResultChan
}
