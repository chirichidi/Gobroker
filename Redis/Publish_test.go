package Redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_Publish(t *testing.T) {
	command := Publish{
		channel:    "topic",
		message:    "message",
		ResultChan: make(chan int64, 1),
	}
	redis := InitRedis(RedisMock())
	command.cmd(redis)

	assert.Equal(t, int64(0), <-command.ResultChan) // 구독자가 없으므로 메시지 받은 수 0
}
