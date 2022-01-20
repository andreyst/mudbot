package botutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasLinePrefix(t *testing.T) {
	assert.Equal(t, HasLinePrefix("test\nmoretest", "moretest"), true)
	assert.Equal(t, HasLinePrefix("test\nmoretest", "test"), true)
	assert.Equal(t, HasLinePrefix("test\nmoretest", "missing"), false)
}

func TestHasLinePrefixes(t *testing.T) {
	assert.Equal(t, HasAnyLinePrefix("test\nmoretest", []string{"no", "moretest"}), true)
	assert.Equal(t, HasAnyLinePrefix("test\nmoretest", []string{"test", "no"}), true)
	assert.Equal(t, HasAnyLinePrefix("test\nmoretest", []string{"no", "moreno"}), false)
}
