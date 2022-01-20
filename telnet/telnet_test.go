package telnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripTelnet(t *testing.T) {
	assert.Equal(t, []byte("test"), StripTelnet([]byte("test")))
	assert.Equal(t, []byte{'x'}, StripTelnet([]byte{IAC, COMMAND_IAC, 'x'}))
	assert.Equal(t, []byte{'x'}, StripTelnet([]byte{IAC, COMMAND_SB, OPTION_MCCPV2, IAC, COMMAND_SE, 'x'}))
	assert.Equal(t, []byte{}, StripTelnet([]byte{}))
	assert.Equal(t, []byte{'x'}, StripTelnet([]byte{'x', IAC}))
	// TODO: Should fail with not a command instead
	assert.Equal(t, []byte{}, StripTelnet([]byte{IAC, 'x'}))
	assert.Equal(t, []byte{}, StripTelnet([]byte{IAC}))
}
