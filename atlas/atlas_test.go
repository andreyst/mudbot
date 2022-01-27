package atlas

import "testing"

func TestAtlas_RecordRoomDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas(false)
	a.RecordRoom(Room{})
}

func TestAtlas_RecordCannotMoveFeedbackDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas(false)
	a.RecordCannotMoveFeedback()
}
