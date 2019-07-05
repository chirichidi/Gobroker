package Redis

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redis_Subscribe(t *testing.T) {

	// 1. subscribe
	r := RedisMock()
	commandSub := Subscribe{
		channel:    "topic",
		ResultChan: make(chan *redis.PubSub, 1),
	}
	commandSub.cmd(InitRedis(r))
	time.Sleep(10 * time.Millisecond)

	// 2. publish
	commandPub := Publish{
		channel: "topic",
		message: "message",
	}
	commandPub.cmd(InitRedis(RedisMock()))

	// then
	assert.Equal(t, int64(1), <-commandSub.ResultChan) // 구독하는거 하나 있는 것과 비교
}
