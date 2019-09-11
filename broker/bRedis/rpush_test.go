package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_RPush(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	command := RPush{
		Key:        "Key",
		Value:      "Value",
		ResultChan: make(chan broker.CmdResult, 1),
	}

	//when
	command.Cmd(InitBRedis(r))

	//then
	assert.Equal(t, "Value", r.BLPop(0, "Key").Val()[1])
}
