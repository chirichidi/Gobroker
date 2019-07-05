package Redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_BRPop(t *testing.T) {
	//given
	r := RedisMock()
	key1 := "key1"
	value1 := "value1"
	r.RPush(key1, value1)

	//when
	command := BRPop{
		time:       0,
		key:        "key1",
		ResultChan: make(chan string, 10),
	}
	redis := InitRedis(r)
	command.cmd(redis)

	//then
	assert.Equal(t, "value1", <-command.ResultChan)
}
