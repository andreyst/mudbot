package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_LinkRooms(t *testing.T) {
	a := NewAtlas()
	a.RecordRoom(NewPrefilledRoom("Room 1", "", []Direction{}, []string{}, []string{}))
	a.RecordMovement(DIRECTION_WEST)
	a.RecordRoom(NewPrefilledRoom("Room 2", "", []Direction{}, []string{}, []string{}))
	a.RecordMovement(DIRECTION_WEST)
	a.RecordRoom(NewPrefilledRoom("Room 3", "", []Direction{}, []string{}, []string{}))

	a.LinkRooms(1, DIRECTION_EAST, 3, DIRECTION_WEST)

	assert.Equal(t, 3, a.Rooms[1].Exits[DIRECTION_EAST])
	assert.Equal(t, 1, a.Rooms[3].Exits[DIRECTION_WEST])
}
