package atlas

import (
	"errors"
	"fmt"
)

type Coordinates struct {
	X int64
	Y int64
	Z int64
}

func (c *Coordinates) AddDir(dir Direction) {
	switch dir {
	case DIRECTION_NORTH:
		c.Y -= 1
	case DIRECTION_SOUTH:
		c.Y += 1
	case DIRECTION_WEST:
		c.X -= 1
	case DIRECTION_EAST:
		c.X += 1
	case DIRECTION_UP:
		c.Z += 1
	case DIRECTION_DOWN:
		c.Z -= 1
	default:
		panic(errors.New("do not know how to subtract direction from coords: %v" + dir.String()))
	}
}

func (c *Coordinates) Shift(dir Direction) {
	switch dir {
	case DIRECTION_NORTH:
		c.Y -= 1
	case DIRECTION_SOUTH:
		c.Y += 1
	case DIRECTION_WEST:
		c.X -= 1
	case DIRECTION_EAST:
		c.X += 1
	case DIRECTION_UP:
		c.Z += 1
	case DIRECTION_DOWN:
		c.Z -= 1
	default:
		panic(errors.New(fmt.Sprintf("Unknown direction to shift room in: %v", dir)))
	}
}
