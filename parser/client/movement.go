package client

import "mudbot/atlas"

var north = "north"
var south = "south"
var west = "west"
var east = "east"
var up = "up"
var down = "down"

func (p *Parser) ParseMovement(s string) (dir atlas.Direction, hasMoved bool) {
	hasMoved = true

	if s == north {
		dir = atlas.DIRECTION_NORTH
	} else if s == south {
		dir = atlas.DIRECTION_SOUTH
	} else if s == west {
		dir = atlas.DIRECTION_WEST
	} else if s == east {
		dir = atlas.DIRECTION_EAST
	} else if s == up {
		dir = atlas.DIRECTION_UP
	} else if s == down {
		dir = atlas.DIRECTION_DOWN
	} else {
		hasMoved = false
	}

	return
}
