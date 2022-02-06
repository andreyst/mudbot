package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_DeleteRoom(t *testing.T) {
	a := NewAtlas()

	room1 := a.AddRoomWithoutMovement(NewRoom("Room1"))
	room2 := a.AddRoomWithoutMovement(NewRoom("Room2"))
	room3 := a.AddRoomWithoutMovement(NewRoom("Room3"))

	a.LinkRooms(room1, DIRECTION_WEST, room2, DIRECTION_EAST)
	a.LinkRooms(room2, DIRECTION_WEST, room3, DIRECTION_EAST)

	a.DeleteRoom(room2)

	assert.Equal(t, 2, len(a.Rooms))
	assert.Equal(t, int64(0), room3.Exits[DIRECTION_EAST])
	assert.Equal(t, int64(0), room1.Exits[DIRECTION_WEST])
}
