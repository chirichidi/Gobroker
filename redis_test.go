package Broker

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setUp() *_redis {
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	})

	r := InitRedis(client)
	return r
}

func _pubSubChannels(r *_redis) *redis.StringSliceCmd {
	return r.Client.PubSubChannels("*")
}

func TestRegister(t *testing.T) {
	r := setUp()
	r.Ping()
}

func Test_redis_RPush(t *testing.T) {
	r := setUp()
	r.RPush("test key", "test value")
}

func Test_redis_BRPop(t *testing.T) {
	r := setUp()
	key1 := "key1"
	value1 := "value1"
	r.RPush(key1, value1)

	key2 := "key2"
	value2 := "value2"
	r.RPush(key2, value2)

	key3 := "key3"
	value3 := "value3"
	r.RPush(key3, value3)

	result := r.BRPop(0, key1, key2)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, string(result[0]), value1)
	assert.Equal(t, string(result[1]), value2)
}

func Test_redis_Publish(t *testing.T) {
	r := setUp()
	n := r.Publish("topic", "message")

	assert.Equal(t, int64(0), n) // 구독자가 없으므로 메시지 받은 수 0
}

func Test_redis_Subscribe(t *testing.T) {
	r := setUp()
	r.Subscribe("topic")
	time.Sleep(1*time.Second)

	n := r.Publish("topic", "message")
	assert.Equal(t, int64(1), n)
}

func Test_redis_ReceiveMessage(t *testing.T) {
	r := setUp()
	ps := r.Subscribe("topic2")
	time.Sleep(1*time.Second)

	n := r.Publish("topic2", "message")
	assert.Equal(t, int64(1), n)

	cnt := 0
	onDoneChan := make(chan struct{})

	for {
		if cnt >= 1 {
			break
		}

		m := r.ReceiveMessage(ps[0], onDoneChan)
		assert.Equal(t, "message", m.Payload)
		cnt++
	}
}

func Test_redis_ReceiveMessage2(t *testing.T) { // 종료 테스트
	r := setUp()
	ps := r.Subscribe("topic3")
	time.Sleep(1*time.Second)

	n := r.Publish("topic3", "message")
	assert.Equal(t, int64(1), n)

	onDoneChan := make(chan struct{})

	go func() {
		r.ReceiveMessage(ps[0], onDoneChan)
	}()

	onDoneChan <- struct{}{}
	assert.Equal(t, 0,len(onDoneChan))
}