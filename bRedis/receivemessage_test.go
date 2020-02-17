package bRedis

import (
	"dancechanlibrary/broker"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func Test_redis_ReceiveMessage(t *testing.T) {
	//given
	r := RedisMock()
	defer r.Close()
	ps := r.Subscribe("topic")
	defer ps.Close()
	time.Sleep(10 * time.Millisecond)

	command := ReceiveMessage{
		PubSub:     ps,
		ResultChan: make(chan broker.CmdResult, 10),
	}

	//when
	command.Cmd(InitBRedis(r))
	time.Sleep(1 * time.Millisecond)

	//then
	isDone := false

	for i := 0; i < 5; i++ {
		r.Publish("topic", i)

	LOOP:
		for {
			assert.Equal(t, strconv.Itoa(i), command.Result().Result)
			isDone = true
			break LOOP
		}
	}

	assert.Equal(t, true, isDone)
}

func Test_redis_ReceiveMessage2(t *testing.T) { // 종료 테스트
	//given
	r := RedisMock()
	defer r.Close()
	ps := r.Subscribe("topic3")
	defer ps.Close()
	time.Sleep(10 * time.Millisecond)

	command := ReceiveMessage{
		PubSub:     ps,
		ResultChan: make(chan broker.CmdResult, 10),
	}

	//when
	command.Cmd(InitBRedis(r))
	time.Sleep(1 * time.Millisecond)

	command.Stop()
	time.Sleep(1 * time.Millisecond)

	//then
	assert.Equal(t, 0, len(command.onDoneChan))
}
