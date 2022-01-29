package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoom_Shift(t *testing.T) {
	r := NewRoom()
	r.Shift(DIRECTION_NORTH)

	assert.Equal(t, int64(0), r.Coordinates.X)
	assert.Equal(t, int64(-1), r.Coordinates.Y)
	assert.Equal(t, int64(0), r.Coordinates.Z)
}
