package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_BLPop(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	key1 := "key1"
	value1 := "value1"
	r.RPush(key1, value1)

	//when
	command := BLPop{
		Time:       0,
		Key:        "key1",
		ResultChan: make(chan broker.CmdResult, 10),
	}
	redis := InitBRedis(r)
	command.Cmd(redis)

	//then
	assert.Equal(t, "value1", command.Result().Result)
}
