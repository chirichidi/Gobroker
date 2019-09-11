package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_Publish(t *testing.T) {
	//given
	command := Publish{
		Channel:    "topic",
		Message:    "Message",
		ResultChan: make(chan broker.CmdResult, 1),
	}
	redis := InitBRedis(RedisMock())
	defer redis.Clear()

	//when
	command.Cmd(redis)

	//then
	assert.Equal(t, int64(0), command.Result().Result) // 구독자가 없으므로 메시지 받은 수 0
}
