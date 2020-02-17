package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type Publish struct {
	Channel    string
	Message    string
	resultChan chan *redis.IntCmd
}

func (p *Publish) Init() {
	p.resultChan = make(chan *redis.IntCmd, 1)
}

func (p *Publish) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(p.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.Publish(p.Channel, p.Message)

	select {
	case p.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (p *Publish) Result() (int64, error) {
	for o := range p.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return 0, nil
}
