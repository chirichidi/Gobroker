package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type Get struct {
	Key        string
	resultChan chan *redis.StringCmd `json:"-"`
}

func (g *Get) Init() {
	g.resultChan = make(chan *redis.StringCmd, 1)
}

func (g *Get) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(g.resultChan)
	}()

	r := worker.(*bRedis)
	cmd := r.client.Get(g.Key)

	select {
	case g.resultChan <- cmd:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (g *Get) Result() (string, error) {
	for o := range g.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return "", nil
}
