package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_UnlinkRooms(t *testing.T) {
	a := NewAtlas()

	room1 := a.AddRoomWithoutMovement(NewRoom("Room1"))
	room2 := a.AddRoomWithoutMovement(NewRoom("Room2"))
	room3 := a.AddRoomWithoutMovement(NewRoom("Room3"))

	a.LinkRooms(room1, DIRECTION_WEST, room2, DIRECTION_EAST)
	a.LinkRooms(room1, DIRECTION_NORTH, room2, DIRECTION_SOUTH)
	a.LinkRooms(room1, DIRECTION_UP, room2, DIRECTION_DOWN)

	a.LinkRooms(room1, DIRECTION_DOWN, room3, DIRECTION_UP)
	a.LinkRooms(room2, DIRECTION_UP, room3, DIRECTION_DOWN)

	a.UnlinkRooms(room1, room2)

	exitsCount := 0
	for _, exitRoomId := range room1.Exits {
		assert.Contains(t, []int64{0, room3.Id}, exitRoomId)
		exitsCount++
	}
	for _, exitRoomId := range room2.Exits {
		assert.Contains(t, []int64{0, room3.Id}, exitRoomId)
		exitsCount++
	}

	assert.Greater(t, exitsCount, 0)
}
