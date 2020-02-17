package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Set(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	command := Set{
		Key:   "SetTest",
		Value: "SetValue",
	}
	redis := InitBRedis(r)

	//when
	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Nil(t, err, nil)
	assert.Equal(t, "OK", result)

	getResult, err := r.Get("SetTest").Result()
	assert.Equal(t, "SetValue", getResult)
	assert.Nil(t, err)
}
