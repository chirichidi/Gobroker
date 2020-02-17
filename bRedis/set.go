package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/go-redis/redis"
	"time"
)

type Set struct {
	Key        string
	Value      interface{}
	Time       time.Duration
	resultChan chan *redis.StatusCmd
}

func (s *Set) Init() {
	s.resultChan = make(chan *redis.StatusCmd, 1)
}

func (s *Set) Cmd(worker broker.Worker) broker.ErrorCodeType {
	defer func() {
		close(s.resultChan)
	}()

	r := worker.(*bRedis)
	result := r.client.Set(s.Key, s.Value, s.Time)

	select {
	case s.resultChan <- result:
		{
		}
	default:
		return broker.BROKER_ERROR_CODE_CUSTOM_CHANNEL_OVERFLOW
	}

	return broker.BROKER_ERROR_CODE_SUCCESS
}

func (s *Set) Result() (string, error) {
	for o := range s.resultChan { // 채널 닫힐 때까지 blocking
		return o.Result()
	}
	return "", nil
}
