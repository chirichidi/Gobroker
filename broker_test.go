package Broker

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func brokerMock() *Manager {
	m := new(Manager)
	m.Init(10).Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	}))

	return m
}

func TestManager_Register(t *testing.T) {
	m := brokerMock()
	assert.Equal(t, "PONG", m._redis.Ping())
}

func TestManager_Up(t *testing.T) {
	m := brokerMock()
	m.Up()

	// up() 동작이 잘 수행중인지 확인하기 위해서 RPush() 보내봄
	param := RedisParam{
		command: COMMAND_REDIS_RPUSH,
		RedisParamDetails: RedisParamDetails{
			Key:    "key",
			Values: []string{"value"},
		},
	}
	m.pushCommand(&param)
	time.Sleep(100 * time.Microsecond)

	param2 := RedisParam{
		command:    COMMAND_REDIS_BRPOP,
		ResultChan: make(chan interface{}, 1),
		RedisParamDetails: RedisParamDetails{
			Keys: []string{"key"},
		},
	}
	m.pushCommand(&param2)
	result := <-param2.ResultChan //result parsing
	s := result.([]string)
	assert.Equal(t, "value", s[0])
}

func TestManager_Stop(t *testing.T) {
	m := brokerMock()
	m.Up()
	time.Sleep(100 * time.Microsecond)

	m.Stop()
	assert.Equal(t, int32(0), m._CommandChanSize)
}

//TODO up 을 여러번 하면? 혹은 Register 를 여러번 하면??
