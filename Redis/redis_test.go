package Redis

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func RedisMock() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     5,
	})

	return client
}

func TestRegister(t *testing.T) {
	r := RedisMock()
	r.Ping()
}
