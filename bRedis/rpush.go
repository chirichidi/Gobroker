package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type RPush struct {
	Key        string
	Value      string
	resultChan chan *redis.IntCmd
}

func (rp *RPush) Init() {
	rp.resultChan = make(chan *redis.IntCmd, 1)
}

func (rp *RPush) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(rp.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.RPush(rp.Key, rp.Value)

	select {
	case rp.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (rp *RPush) Result() (int64, error) {
	for o := range rp.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return 0, nil
}
