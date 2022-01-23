package atlas

import (
	"errors"
	"fmt"
)

func (d Direction) Opposite() Direction {
	switch d {
	case DIRECTION_NORTH:
		return DIRECTION_SOUTH
	case DIRECTION_SOUTH:
		return DIRECTION_NORTH
	case DIRECTION_WEST:
		return DIRECTION_EAST
	case DIRECTION_EAST:
		return DIRECTION_WEST
	case DIRECTION_UP:
		return DIRECTION_DOWN
	case DIRECTION_DOWN:
		return DIRECTION_UP
	default:
		panic(errors.New(fmt.Sprintf("Unknown opposite for direction: %s", d)))
	}
}
