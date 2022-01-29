package client

import (
	"github.com/oriser/regroup"
	"mudbot/atlas"
	"strconv"
)

var shiftMatcher = regroup.MustCompile(`/shift (?P<RoomId>\d+) (?P<Direction>[NSWEUD])`)

func (p *Parser) ParseShift(s string) (roomId int64, direction atlas.Direction, ok bool) {
	match, _ := shiftMatcher.Groups(s)
	if match == nil {
		return
	}

	ok = true
	roomId, _ = strconv.ParseInt(match["RoomId"], 10, 64)
	direction = atlas.NewDirection(match["Direction"])

	return
}
