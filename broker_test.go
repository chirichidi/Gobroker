package Broker

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func redisUp() *Manager {
	m := new(Manager)

	m.Init(3).Register(redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	})).Up()

	return m
}

func TestManager_Register(t *testing.T) {
	m := redisUp()
	assert.Equal(t, "PONG", m._redis.Ping()) //TODO 왜 여기서 _redis 의 클라이언트가 nil 이지???
}

func TestManager_Stop(t *testing.T) {

}
