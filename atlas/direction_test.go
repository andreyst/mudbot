package atlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirection_Opposite(t *testing.T) {
	assert.Equal(t, DIRECTION_DOWN, DIRECTION_UP.Opposite())
}
