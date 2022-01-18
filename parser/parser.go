package parser

import (
	"bytes"
	"mudbot/botutil"
	"mudbot/char"

	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.SugaredLogger

	char *char.Char
}

func NewParser() *Parser {
	parser := Parser{
		logger: botutil.NewLogger("parser"),
		char:   char.NewChar(),
	}

	return &parser

}
func (p *Parser) Parse(chunk []byte) {
	chunk = bytes.ReplaceAll(chunk, []byte{'\r', '\n'}, []byte{'\n'})
	p.char.ParseScore(string(chunk))
}
