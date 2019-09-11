package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redis_Subscribe(t *testing.T) {
	// 1. subscribe
	r := RedisMock()
	defer r.Close()
	commandSub := Subscribe{
		Channel:    "topic",
		ResultChan: make(chan broker.CmdResult, 3),
	}
	commandSub.Cmd(InitBRedis(r))
	time.Sleep(1 * time.Millisecond)

	// 2. publish
	commandPub := Publish{
		Channel:    "topic",
		Message:    "Message",
		ResultChan: make(chan broker.CmdResult, 3),
	}
	commandPub.Cmd(InitBRedis(RedisMock()))
	time.Sleep(1 * time.Millisecond)

	// then
	assert.Equal(t, int64(1), commandPub.Result().Result) // 구독하는거 하나 있는 것과 비교
}
