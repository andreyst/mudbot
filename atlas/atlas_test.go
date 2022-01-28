package atlas

import "testing"

func TestAtlas_RecordRoomDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas(false)
	r := Room{}
	a.RecordRoom(&r)
}

func TestAtlas_RecordCannotMoveFeedbackDoesNotCrashWithoutMovement(t *testing.T) {
	a := NewAtlas(false)
	a.RecordCannotMoveFeedback()
}
