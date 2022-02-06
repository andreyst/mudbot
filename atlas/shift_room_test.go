package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_Shift(t *testing.T) {
	atlas := NewAtlas()

	var room *Room

	room = NewFilledRoom("Room 1", "", []Direction{DIRECTION_WEST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_WEST)
	room = NewFilledRoom("Room 2", "", []Direction{DIRECTION_WEST, DIRECTION_EAST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_WEST)
	room = NewFilledRoom("Room 3", "", []Direction{DIRECTION_NORTH, DIRECTION_WEST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_NORTH)
	room = NewFilledRoom("Room 4", "", []Direction{DIRECTION_NORTH, DIRECTION_SOUTH}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_NORTH)
	room = NewFilledRoom("Room 5", "", []Direction{DIRECTION_NORTH, DIRECTION_EAST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_EAST)
	room = NewFilledRoom("Room 6", "", []Direction{DIRECTION_WEST, DIRECTION_EAST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_EAST)
	room = NewFilledRoom("Room 7", "", []Direction{DIRECTION_WEST, DIRECTION_EAST}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_EAST)
	room = NewFilledRoom("Room 8", "", []Direction{DIRECTION_WEST, DIRECTION_SOUTH}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_SOUTH)
	room = NewFilledRoom("Room 9", "", []Direction{DIRECTION_WEST, DIRECTION_NORTH}, []string{}, []string{})
	atlas.RecordRoom(room)

	atlas.RecordMovement(DIRECTION_WEST)
	room = NewFilledRoom("Room 10", "", []Direction{DIRECTION_EAST}, []string{}, []string{})
	atlas.RecordRoom(room)

	yBefore := make(map[int64]int64)
	for i := 1; i < len(atlas.Rooms)+1; i++ {
		yBefore[int64(i)] = atlas.Rooms[int64(i)].Coordinates.Y
	}

	atlas.ShiftRoom(2, DIRECTION_NORTH)

	for i := 1; i < len(atlas.Rooms)+1; i++ {
		assert.Equal(t, yBefore[int64(i)]-1, atlas.Rooms[int64(i)].Coordinates.Y)
	}
}
