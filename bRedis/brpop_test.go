package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_BRPop(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	key1 := "key1"
	value1 := "value1"
	r.RPush(key1, value1)

	//when
	command := BRPop{
		Time: 0,
		Keys: []string{"key1"},
	}
	redis := InitBRedis(r)
	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Nil(t, err)
	assert.Equal(t, "key1", result[0])
	assert.Equal(t, "value1", result[1])
}
