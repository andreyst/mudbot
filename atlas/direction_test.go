package atlas

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirection_String(t *testing.T) {
	assert.Equal(t, "N", DIRECTION_NORTH.String())
}

func TestDirection_Opposite(t *testing.T) {
	assert.Equal(t, DIRECTION_DOWN, DIRECTION_UP.Opposite())
}

func TestDirection_MarshalText(t *testing.T) {
	res, marshalErr := json.Marshal(DIRECTION_UP)
	assert.Nil(t, marshalErr)
	assert.Equal(t, []byte(`"U"`), res)
}

func TestDirection_UnmarshalText(t *testing.T) {
	var dir Direction
	unmarshalErr := json.Unmarshal([]byte(`"U"`), &dir)
	assert.Nil(t, unmarshalErr)
	assert.Equal(t, DIRECTION_UP, dir)
}
