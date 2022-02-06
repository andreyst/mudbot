package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetFirstRoom(t *testing.T) {
	rooms := make(Rooms)
	assert.Nil(t, getFirstRoom(rooms))
	room1 := NewRoom("room1")
	room2 := NewRoom("room2")
	rooms[1] = room1
	rooms[2] = room2
	assert.NotNil(t, getFirstRoom(rooms))
}
