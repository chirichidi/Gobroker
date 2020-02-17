package bRedis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func RedisMock() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "12villeDance@",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	})

	client.FlushAll()
	return client
}

func TestBRedis_Ping(t *testing.T) {
	//given
	c := RedisMock()
	defer c.Close()

	//when
	r := InitBRedis(c)
	b := r.Ping()

	//then
	assert.Equal(t, true, b)
}
