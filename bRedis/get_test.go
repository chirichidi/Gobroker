package bRedis

import (
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
		Key: "GetTest",
	}
	redis := InitBRedis(r)

	//when

	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Nil(t, err)
	assert.Equal(t, "GetValue", result)
}
