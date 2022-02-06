package atlas

import (
	"github.com/stretchr/testify/assert"
	"mudbot/atlas/server"
	"testing"
)

func TestAtlas_LinkRoom(t *testing.T) {
	a := NewAtlas()
	room1 := a.AddRoomWithoutMovement(NewRoom("Room 1"))
	room2 := a.AddRoomWithoutMovement(NewRoom("Room 2"))

	a.LinkRoom(room1, DIRECTION_EAST, room2)

	assert.Equal(t, room2.Id, a.Rooms[room1.Id].Exits[DIRECTION_EAST])
	assert.Equal(t, 0, len(room2.Exits))
}

func TestAtlas_OnLinkRoom(t *testing.T) {
	a := NewAtlas()
	room1 := a.AddRoomWithoutMovement(NewRoom("Room 1"))
	assert.NotPanics(t, func() {
		a.onLinkRoom(server.LinkRoomCommand{FromRoomId: int(room1.Id), FromRoomExit: "S", ToRoomId: 999})
	})
}
