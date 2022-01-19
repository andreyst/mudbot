package botutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripAnsi1(t *testing.T) {
	ref := "asdlkalds"
	res := StripAnsi(ref)

	assert.Equal(t, ref, res)
}

func TestStripAnsi2(t *testing.T) {
	in := "[0;32m36ж[0;0m [0;32m100б[0;0m 2952о 30м Выходы:Ю>"
	ref := "36ж 100б 2952о 30м Выходы:Ю>"
	res := StripAnsi(in)

	assert.Equal(t, ref, res)
}
