package bRedis

import (
	"dancechanlibrary/broker"
	"time"
)

type BRPop struct {
	Time       time.Duration
	Key        string                // TODO 성능때문에 single parameter 사용
	ResultChan chan broker.CmdResult `json:"-"`
}

func (bp *BRPop) Cmd(worker broker.Worker) broker.ErrorCodeType {
	r := worker.(*bRedis)
	result, err := r.client.BRPop(bp.Time, bp.Key).Result()

	select {
	case bp.ResultChan <- broker.CmdResult{Result: result[1], Err: err}:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (bp *BRPop) Result() broker.CmdResult {
	return <-bp.ResultChan
}
