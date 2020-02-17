package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
	"time"
)

type BRPop struct {
	Time       time.Duration
	Keys       []string
	resultChan chan *redis.StringSliceCmd
}

func (bp *BRPop) Init() {
	bp.resultChan = make(chan *redis.StringSliceCmd, len(bp.Keys))
}

func (bp *BRPop) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(bp.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.BRPop(bp.Time, bp.Keys...)

	select {
	case bp.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (bp *BRPop) Result() ([]string, error) {
	for o := range bp.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return []string{}, nil
}
