package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
)

type Subscribe struct {
	Channel    string             //TODO 성능 이슈로 한개만
	resultChan chan *redis.PubSub `json:"-"`
}

func (s *Subscribe) Init() {
	s.resultChan = make(chan *redis.PubSub, 1)
}

func (s *Subscribe) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(s.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.Subscribe(s.Channel)

	select {
	case s.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (s *Subscribe) Result() *redis.PubSub {
	for o := range s.resultChan { // 채널 닫힐 때까지 blocking
		return o
	}
	return nil
}
