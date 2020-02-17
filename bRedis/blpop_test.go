package bRedis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		Time: 0,
		Keys: []string{"key1"},
	}
	redis := InitBRedis(r)
	command.Cmd(redis)

	//then
	result, err := command.Result()
	assert.Equal(t, err, nil)
	assert.Equal(t, "key1", result[0])
	assert.Equal(t, "value1", result[1])
}

func Test_result_메소드_동작_방식_확인(t *testing.T) {
	//given
	a := make(chan int32)
	go func() {
		time.Sleep(time.Second * 3)
		a <- 1
		close(a)
		a = nil
	}()

	//when
	for key := range a {
		assert.Equal(t, key, int32(1))
	}

	//then
	assert.Nil(t, a)
}
