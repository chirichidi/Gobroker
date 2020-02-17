package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_RPush(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	command := RPush{
		Key:   "Key",
		Value: "Value",
	}

	//when
	command.Cmd(InitBRedis(r))

	//then
	assert.Equal(t, "Value", r.BLPop(0, "Key").Val()[1])
}
