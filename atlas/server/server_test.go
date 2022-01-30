package server

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestNewServer(t *testing.T) {
	{
		msg := Message{
			Command: "wow",
			ShiftCommand: &ShiftCommand{
				RoomId:    111,
				Direction: "S",
			},
		}
		msgBytes, _ := json.Marshal(msg)
		t.Logf("Message: %+v", string(msgBytes))
	}

	{
		var msg Message
		s := `{"Command":"wow"}}`
		json.Unmarshal([]byte(s), &msg)
		t.Logf(spew.Sprintf("Unmarshalled: %+v", msg))
	}
}
