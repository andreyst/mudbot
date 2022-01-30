package client

import (
	"github.com/oriser/regroup"
	"mudbot/atlas"
	"strconv"
)

var shiftMatcher = regroup.MustCompile(`/map shift (?P<RoomId>\d+) (?P<Direction>[NSWEUD])`)

func (p *Parser) ParseShift(s string) (roomId int64, direction atlas.Direction, ok bool) {
	match, _ := shiftMatcher.Groups(s)
	if match == nil {
		return
	}

	ok = true
	var parseIntErr error
	roomId, parseIntErr = strconv.ParseInt(match["RoomId"], 10, 64)
	if parseIntErr != nil {
		p.logger.Debugf("Cannot parse int: %v", match["RoomId"])
		ok = false
		return
	}

	direction, ok = atlas.NewDirection(match["Direction"])

	return
}
