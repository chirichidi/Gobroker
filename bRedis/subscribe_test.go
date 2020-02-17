package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redis_Subscribe(t *testing.T) {
	// 1. subscribe
	r := RedisMock()
	defer r.Close()
	commandSub := Subscribe{
		Channel: "topic",
	}
	commandSub.Cmd(InitBRedis(r))
	time.Sleep(1 * time.Millisecond)

	// 2. publish
	commandPub := Publish{
		Channel: "topic",
		Message: "Message",
	}
	commandPub.Cmd(InitBRedis(RedisMock()))
	time.Sleep(1 * time.Millisecond)

	// then
	result, err := commandPub.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result) // 구독하는거 하나 있는 것과 비교
}
