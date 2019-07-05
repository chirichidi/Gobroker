package Redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_redis_RPush(t *testing.T) {
	r := RedisMock()
	command := RPush{
		key:   "key",
		value: "value",
	}
	result := command.cmd(InitRedis(r))

	assert.Equal(t, true, result)
}
