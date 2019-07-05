package Redis

import (
	"Broker"
	"github.com/go-redis/redis"
)

type _redis struct {
	Client *redis.Client
}

//TODO 브로커 생성시점에 같이 생성하는게 가독성상 좋다
func InitRedis(client *redis.Client) *_redis {
	defer Broker.PrintPanicStack()

	// 0. check
	_, err := client.Ping().Result()
	if err != nil {
		panic(err.Error())
		return nil
	}

	// 1. register
	_client := &_redis{
		Client: client,
	}
	return _client
}

func (r *_redis) _Clear() {
	r.Client = nil
}

func (r *_redis) Ping() bool {
	defer Broker.PrintPanicStack()

	p := r.Client.Ping()
	if p.Err() != nil {
		panic(p.Err())
	}

	return true
}

func (r *_redis) PushChan(ch chan Broker.WorkerCommand, c Broker.Commander) {
	defer Broker.PrintPanicStack()

	ch <- Broker.WorkerCommand{
		Worker:    r,
		Commander: c,
	}
}
