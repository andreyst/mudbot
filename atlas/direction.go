package atlas

import (
	"errors"
	"fmt"
	"strings"
)

type Direction int

const (
	DIRECTION_UNKNOWN Direction = iota
	DIRECTION_NORTH
	DIRECTION_SOUTH
	DIRECTION_WEST
	DIRECTION_EAST
	DIRECTION_UP
	DIRECTION_DOWN
)

func NewDirection(s string) (dir Direction, ok bool) {
	ok = true

	switch strings.ToUpper(s) {
	case "N":
		dir = DIRECTION_NORTH
	case "S":
		dir = DIRECTION_SOUTH
	case "W":
		dir = DIRECTION_WEST
	case "E":
		dir = DIRECTION_EAST
	case "U":
		dir = DIRECTION_UP
	case "D":
		dir = DIRECTION_DOWN
	default:
		ok = false
	}

	return
}

func (d Direction) String() string {
	switch d {
	case DIRECTION_NORTH:
		return "N"
	case DIRECTION_SOUTH:
		return "S"
	case DIRECTION_WEST:
		return "W"
	case DIRECTION_EAST:
		return "E"
	case DIRECTION_UP:
		return "U"
	case DIRECTION_DOWN:
		return "D"
	default:
		panic(errors.New(fmt.Sprintf("unknown direction to str: %d", d)))
	}
}

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
		panic(errors.New(fmt.Sprintf("unknown opposite for direction: %s", d)))
	}
}

func (d *Direction) UnmarshalText(bytes []byte) error {
	s := string(bytes)
	dir, ok := NewDirection(s)
	if !ok {
		return errors.New(fmt.Sprintf("Unknown direction %v", s))

	}

	*d = dir
	return nil
}

func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
