package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MGet(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	key1 := "MGetTest1"
	value1 := "MGetValue1"
	r.Set(key1, value1, 0)

	key2 := "MGetTest2"
	value2 := "MGetValue2"
	r.Set(key2, value2, 0)

	command := MGet{
		Keys: []string{"MGetTest1", "MGetTest2", "MGetTest3"}, //MGetTest3 는 없는 녀석
	}
	redis := InitBRedis(r)

	//when

	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Nil(t, err)
	assert.Equal(t, "MGetValue1", result[0])
	assert.Equal(t, "MGetValue2", result[1])
	assert.Equal(t, nil, result[2])
}
