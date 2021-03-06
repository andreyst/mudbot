package atlas

import (
	"github.com/stretchr/testify/assert"
	"mudbot/atlas/server"
	"testing"
)

func TestAtlas_LinkRooms(t *testing.T) {
	a := NewAtlas()
	room1 := a.AddRoomWithoutMovement(NewRoom("Room 1"))
	room2 := a.AddRoomWithoutMovement(NewRoom("Room 2"))

	a.LinkRooms(room1, DIRECTION_EAST, room2, DIRECTION_WEST)

	assert.Equal(t, room2.Id, a.Rooms[room1.Id].Exits[DIRECTION_EAST])
	assert.Equal(t, room1.Id, a.Rooms[room2.Id].Exits[DIRECTION_WEST])
}

func TestAtlas_OnLinkRooms(t *testing.T) {
	a := NewAtlas()
	room1 := a.AddRoomWithoutMovement(NewRoom("Room 1"))
	assert.NotPanics(t, func() {
		a.onLinkRooms(server.LinkRoomsCommand{FromRoomId: int(room1.Id), FromRoomExit: "S", ToRoomId: 999, ToRoomExit: "N"})
	})
}
