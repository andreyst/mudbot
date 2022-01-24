package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirection_String(t *testing.T) {
	assert.Equal(t, "N", DIRECTION_NORTH.String())
}

func TestDirection_Opposite(t *testing.T) {
	assert.Equal(t, DIRECTION_DOWN, DIRECTION_UP.Opposite())
}
