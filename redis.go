package Broker

import (
	"github.com/go-redis/redis"
	"time"
)

type _redis struct {
	Client *redis.Client
}

type RedisParam struct {
	command    int32
	ResultChan chan interface{} // app 의 결과 수신용 채널
	RedisParamDetails
}

type RedisParamDetails struct {
	Key        string
	Keys       []string
	Value      string
	Values     []string
	Timeout    time.Duration
	Channel    string
	Channels   []string
	Message    string
	PubSub     *redis.PubSub
	OnDoneChan chan struct{}
}

func InitRedis(client *redis.Client) *_redis {
	defer PrintPanicStack()

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

func (r *_redis) _Connect() {

}

func (r *_redis) _Clear() {
	r.Client = nil
}

func (r *_redis) Ping() string {
	defer PrintPanicStack()

	p := r.Client.Ping()
	if p.Err() != nil {
		panic(p.Err())
	}

	return p.Val()
}

func (r *_redis) RPush(key string, values []string) {
	defer PrintPanicStack()

	for _, value := range values { //TODO 흠..............
		_, err := r.Client.RPush(key, value).Result()
		if err != nil {
			panic(err.Error())
			return
		}
	}
}

func (r *_redis) BRPop(timeout time.Duration, keys []string) []string {
	defer PrintPanicStack()

	result := make([]string, len(keys))
	for i, key := range keys {
		brPop := r.Client.BRPop(timeout, key) //TODO 흠.......

		val := brPop.Val()
		result[i] = val[1]
	}
	return result
}

func (r *_redis) BLPop(timeout time.Duration, keys []string) []string {
	defer PrintPanicStack()

	result := make([]string, len(keys))
	for i, key := range keys {
		brPop := r.Client.BLPop(timeout, key) //TODO 흠.......

		val := brPop.Val()
		result[i] = val[1]
	}
	return result
}

func (r *_redis) Publish(channel string, message string) int64 {
	defer PrintPanicStack()

	n, err := r.Client.Publish(channel, message).Result()
	if err != nil {
		panic(err.Error())
		return -1
	}
	return n
}

func (r *_redis) Subscribe(channels []string) []*redis.PubSub {
	defer PrintPanicStack()

	result := make([]*redis.PubSub, len(channels)) //TODO 흠........
	for i, topic := range channels {
		ps := r.Client.Subscribe(topic)
		result[i] = ps
	}

	return result
}

func (r *_redis) UnSubscribe(pubSub *redis.PubSub, channels []string) {
	defer PrintPanicStack()

	for _, topic := range channels {
		pubSub.Unsubscribe(topic)
	}
}

func (r *_redis) ReceiveMessage(pubSub *redis.PubSub, _onDoneChan <-chan struct{}) *redis.Message {
	defer PrintPanicStack()

	select {
	case <-_onDoneChan:
		pubSub.Close()
	default:
		m, err := pubSub.ReceiveMessage()
		if m == nil { // 서버를 강제 종료할 경우 이 케이스에 걸림
			return nil
		}
		if err != nil {
			panic(err.Error())
			return nil
		}
		return m
	}
	return nil
}
