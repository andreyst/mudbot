package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoordinates_AddDir(t *testing.T) {
	c := Coordinates{}

	c.AddDir(DIRECTION_NORTH)
	c.AddDir(DIRECTION_SOUTH)
	c.AddDir(DIRECTION_SOUTH)
	c.AddDir(DIRECTION_EAST)
	c.AddDir(DIRECTION_WEST)
	c.AddDir(DIRECTION_WEST)
	c.AddDir(DIRECTION_UP)
	c.AddDir(DIRECTION_DOWN)
	c.AddDir(DIRECTION_DOWN)

	assert.Equal(t, Coordinates{X: -1, Y: 1, Z: -1}, c)
}
