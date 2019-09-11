package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	key1 := "GetTest"
	value1 := "GetValue"
	r.Set(key1, value1, 0)

	command := Get{
		Key:        "GetTest",
		ResultChan: make(chan broker.CmdResult, 1),
	}
	redis := InitBRedis(r)

	//when

	command.Cmd(redis)

	//then
	result := command.Result()
	assert.Equal(t, "GetValue", result.Result)
	assert.Nil(t, result.Err)
}
