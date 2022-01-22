package botutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasLinePrefix(t *testing.T) {
	assert.Equal(t, HasLinePrefix("test\nmore test", "more test"), true)
	assert.Equal(t, HasLinePrefix("test\nmore test", "test"), true)
	assert.Equal(t, HasLinePrefix("test\nmore test", "missing"), false)
}

func TestHasLinePrefixes(t *testing.T) {
	assert.Equal(t, HasAnyLinePrefix("test\nmore test", []string{"no", "more test"}), true)
	assert.Equal(t, HasAnyLinePrefix("test\nmore test", []string{"test", "no"}), true)
	assert.Equal(t, HasAnyLinePrefix("test\nmore test", []string{"no", "moreno"}), false)
}
