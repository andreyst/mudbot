package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_ParseClientMessageBadCommands(t *testing.T) {
	s := NewServer(nil)
	s.OnShift = func(cmd ShiftCommand) { assert.Fail(t, "Handler called") }

	var buf string

	buf = `{"Command":""}`
	s.parseClientMessage([]byte(buf))

	buf = `{"Command":"shift"}`
	s.parseClientMessage([]byte(buf))
}

func TestServer_ParseClientMessageWithoutHandlers(t *testing.T) {
	s := NewServer(nil)

	var buf string

	buf = `{"Command":"shift", "ShiftCommand":{"RoomId":1, "Direction":"N"}}`
	s.parseClientMessage([]byte(buf))

	buf = `{"Command":"delete", "DeleteCommand":{"RoomId":1}}`
	s.parseClientMessage([]byte(buf))
}

func TestServer_ParseClientMessage(t *testing.T) {
	s := NewServer(nil)

	refCmd := ShiftCommand{
		RoomId:    111,
		Direction: "W",
	}
	msg := Message{
		Command:      "shift",
		ShiftCommand: &refCmd,
	}
	msgBytes, marshalErr := json.Marshal(msg)
	assert.NoError(t, marshalErr)

	var handlerCalled bool
	s.OnShift = func(cmd ShiftCommand) {
		handlerCalled = true
		assert.Equal(t, refCmd.RoomId, cmd.RoomId)
		assert.Equal(t, refCmd.Direction, cmd.Direction)
	}

	s.parseClientMessage(msgBytes)

	assert.True(t, handlerCalled)
}
