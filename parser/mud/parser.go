package mud

import (
	"go.uber.org/zap"
	"mudbot/atlas"
	"mudbot/botutil"
	"strings"
)

type Parser struct {
	atlas *atlas.Atlas

	logger *zap.SugaredLogger
}

func NewParser(atlas *atlas.Atlas) *Parser {
	p := Parser{
		atlas: atlas,

		logger: botutil.NewLogger("mud_parser"),
	}

	return &p
}

func (p *Parser) Parse(bytes []byte) {
	s := strings.TrimRight(string(bytes), "\r\n")

	room, matchedRoom := p.ParseRoom(s)
	if matchedRoom {
		p.atlas.RecordRoom(&room)
	}

	events := p.ParseFeedback(s)
	for _, e := range events {
		if e == EVENT_CANNOT_MOVE_IN_THIS_DIRECTION ||
			e == EVENT_CANNOT_BECAUSE_RESTING ||
			e == EVENT_CANNOT_BECAUSE_SITTING ||
			e == EVENT_CANNOT_BECAUSE_SLEEPING ||
			e == EVENT_CANNOT_BECAUSE_CLOSED {
			p.atlas.RecordCannotMoveFeedback()
		}
	}
}
