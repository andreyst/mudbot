// Package atlas implements functions to work with MUD world map.
//
// Better name for this package would be "map" but it clashes with data structure name.
package atlas

import (
	"go.uber.org/zap"
	"mudbot/botutil"
)

type Atlas struct {
	movements []Direction

	logger *zap.SugaredLogger
}

func NewAtlas() *Atlas {
	a := Atlas{
		logger: botutil.NewLogger("atlas"),
	}

	return &a
}

func (a *Atlas) RecordMovement(dir Direction) {
	a.movements = append(a.movements, dir)
	a.logger.Debugf("Recorded movement: %v", dir)
}

func (a *Atlas) RecordRoom(room Room) {
	fromStr := "not moving"
	if len(a.movements) > 0 {
		from := a.movements[0]
		fromStr = from.Opposite().String()
		a.movements = a.movements[1:]
	}
	a.logger.Debugf("Recorded room: %+v (no info: %t), moved from %v", room.Name, room.NoInfo, fromStr)
}

func (a *Atlas) RecordCannotMoveFeedback() {
	a.movements = a.movements[1:]
	a.logger.Debugf("Recorded cannot move feedback")
}
