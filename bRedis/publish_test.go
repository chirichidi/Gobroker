package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_Publish(t *testing.T) {
	//given
	command := Publish{
		Channel: "topic",
		Message: "Message",
	}
	redis := InitBRedis(RedisMock())
	defer redis.Clear()

	//when
	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), result) // 구독자가 없으므로 메시지 받은 수 0
}
