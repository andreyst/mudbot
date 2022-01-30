package botutil

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, int64(1), Abs(1))
	assert.Equal(t, int64(1), Abs(-1))
	assert.Equal(t, int64(0), Abs(0))
	assert.Equal(t, int64(math.MaxInt64), Abs(-math.MaxInt64))
}
