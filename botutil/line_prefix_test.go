package botutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasLinePrefix(t *testing.T) {
	assert.Equal(t, true, HasLinePrefix("test\nmore test", "more test"))
	assert.Equal(t, true, HasLinePrefix("test\nmore test", "test"))
	assert.Equal(t, false, HasLinePrefix("test\nmore test", "missing"))
}

func TestHasLinePrefixes(t *testing.T) {
	assert.Equal(t, true, HasAnyLinePrefix("test\nmore test", []string{"no", "more test"}))
	assert.Equal(t, true, HasAnyLinePrefix("test\nmore test", []string{"test", "no"}))
	assert.Equal(t, false, HasAnyLinePrefix("test\nmore test", []string{"no", "moreno"}))
}

func TestHasLineSuffix(t *testing.T) {
	assert.Equal(t, true, HasLineSuffix("test\nmore test", "more test"))
	assert.Equal(t, true, HasLineSuffix("test\nmore test", "test"))
	assert.Equal(t, false, HasLineSuffix("test\nmore test", "missing"))
	assert.Equal(t, false, HasLineSuffix("test no\nmore test no", "test"))
}

func TestHasLineSuffixes(t *testing.T) {
	assert.Equal(t, true, HasAnyLineSuffix("test\nmore test", []string{"no", "more test"}))
	assert.Equal(t, true, HasAnyLineSuffix("test\nmore test", []string{"test", "no"}))
	assert.Equal(t, false, HasAnyLineSuffix("test\nmore test", []string{"no", "more"}))
}
