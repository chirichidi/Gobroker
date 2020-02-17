package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type MGet struct {
	Keys       []string
	resultChan chan *redis.SliceCmd
}

func (g *MGet) Init() {
	g.resultChan = make(chan *redis.SliceCmd, len(g.Keys))
}

func (g *MGet) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(g.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.MGet(g.Keys...)

	select {
	case g.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (g *MGet) Result() ([]interface{}, error) {
	for o := range g.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return make([]interface{}, 0), nil
}
