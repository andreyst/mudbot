package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_RecordRoomDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas()
	r := Room{}
	assert.NotPanics(t, func() { a.RecordRoom(&r) })
}

func TestAtlas_RecordCannotMoveFeedbackDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas()
	assert.NotPanics(t, func() { a.RecordCannotMoveFeedback() })
}

func TestAtlas_RecordRoom(t *testing.T) {
	a := NewAtlas()
	// Dereference room template, so we can record it and then use again without ID
	roomTpl := *NewRoomWithExits("Room1", []Direction{DIRECTION_EAST, DIRECTION_SOUTH})
	room := roomTpl
	a.RecordRoom(&room)
	a.RecordMovement(DIRECTION_EAST)
	a.RecordRoom(NewRoomWithExits("Room2", []Direction{DIRECTION_WEST, DIRECTION_SOUTH}))
	a.RecordMovement(DIRECTION_SOUTH)
	a.RecordRoom(NewRoomWithExits("Room3", []Direction{DIRECTION_NORTH, DIRECTION_EAST}))
	a.RecordMovement(DIRECTION_WEST)
	a.RecordRoom(NewRoomWithExits("Room4", []Direction{DIRECTION_NORTH, DIRECTION_EAST}))
	a.RecordMovement(DIRECTION_NORTH)
	room = roomTpl
	a.RecordRoom(&room)

	assert.Equal(t, 4, len(a.Rooms))
}
