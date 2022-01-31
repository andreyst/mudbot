package client

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

		logger: botutil.NewLogger("client_parser"),
	}

	return &p
}

func (p *Parser) Parse(bytes []byte) {
	s := strings.TrimRight(string(bytes), "\r\n")

	dir, hasMoved := p.ParseMovement(s)
	if hasMoved {
		p.atlas.RecordMovement(dir)
	}

	roomId, shiftDir, isShift := p.ParseShift(s)
	if isShift {
		p.atlas.ShiftRoom(roomId, shiftDir)
	}
}
