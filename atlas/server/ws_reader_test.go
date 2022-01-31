package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_ParseClientMessageBadCommands(t *testing.T) {
	s := NewServer(nil)
	s.OnShiftRoom = func(cmd ShiftRoomCommand) { assert.Fail(t, "Handler called") }

	var buf string

	buf = `{"Command":""}`
	s.parseClientMessage([]byte(buf))

	buf = `{"Command":"shift_room"}`
	s.parseClientMessage([]byte(buf))
}

func TestServer_ParseClientMessageWithoutHandlers(t *testing.T) {
	s := NewServer(nil)

	var buf string

	buf = `{"Command":"shift_room", "ShiftRoomCommand":{"RoomId":1, "Direction":"N"}}`
	s.parseClientMessage([]byte(buf))

	buf = `{"Command":"delete_room", "DeleteRoomCommand":{"RoomId":1}}`
	s.parseClientMessage([]byte(buf))
}

func TestServer_ParseClientMessage(t *testing.T) {
	s := NewServer(nil)

	refCmd := ShiftRoomCommand{
		RoomId:    111,
		Direction: "W",
	}
	msg := Message{
		Command:          "shift_room",
		ShiftRoomCommand: &refCmd,
	}
	msgBytes, marshalErr := json.Marshal(msg)
	assert.NoError(t, marshalErr)

	var handlerCalled bool
	s.OnShiftRoom = func(cmd ShiftRoomCommand) {
		handlerCalled = true
		assert.Equal(t, refCmd.RoomId, cmd.RoomId)
		assert.Equal(t, refCmd.Direction, cmd.Direction)
	}

	s.parseClientMessage(msgBytes)

	assert.True(t, handlerCalled)
}
