package Redis

import (
	"Broker"
	"github.com/go-redis/redis"
)

type Subscribe struct {
	channel    string //TODO 성능 이슈로 한개만
	ResultChan chan *redis.PubSub
}

func (s *Subscribe) cmd(worker Broker.Worker) bool {
	defer Broker.PrintPanicStack()
	if len(s.ResultChan) >= cap(s.ResultChan) {
		return false
	}

	r := worker.(*_redis)
	s.ResultChan <- r.Client.Subscribe(s.channel)

	return true
}
