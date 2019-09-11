package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Set(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	command := Set{
		Key:        "SetTest",
		Value:      "SetValue",
		ResultChan: make(chan broker.CmdResult, 1),
	}
	redis := InitBRedis(r)

	//when
	command.Cmd(redis)

	//then
	result := command.Result()
	assert.Equal(t, "OK", result.Result)
	assert.Nil(t, result.Err)

	getResult, err := r.Get("SetTest").Result()
	assert.Equal(t, "SetValue", getResult)
	assert.Nil(t, err)
}
