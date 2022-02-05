package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtlas_DeleteRoom(t *testing.T) {
	a := NewAtlas()
	a.RecordRoom(NewPrefilledRoom("Room1", "", []Direction{DIRECTION_WEST}, []string{}, []string{}))

	a.RecordMovement(DIRECTION_WEST)
	a.RecordRoom(NewPrefilledRoom("Room2", "", []Direction{DIRECTION_WEST, DIRECTION_EAST}, []string{}, []string{}))

	a.RecordMovement(DIRECTION_WEST)
	a.RecordRoom(NewPrefilledRoom("Room3", "", []Direction{DIRECTION_EAST}, []string{}, []string{}))

	a.DeleteRoom(2)

	assert.Equal(t, 2, len(a.Rooms))
	assert.Equal(t, int64(0), a.Rooms[3].Exits[DIRECTION_EAST])
	assert.Equal(t, int64(0), a.Rooms[1].Exits[DIRECTION_WEST])
}
