package client

import (
	"github.com/stretchr/testify/assert"
	"mudbot/atlas"
	"testing"
)

func TestParser_ParseShift(t *testing.T) {
	a := atlas.NewAtlas(false)
	p := NewParser(a)

	{
		roomId, direction, ok := p.ParseShift("/shift 99 S")
		assert.Equal(t, int64(99), roomId)
		assert.Equal(t, atlas.DIRECTION_SOUTH, direction)
		assert.True(t, ok)
	}
	{
		_, _, ok := p.ParseShift("/shift 99 ???")
		assert.False(t, ok)
	}
	{
		_, _, ok := p.ParseShift("no shift")
		assert.False(t, ok)
	}
}
