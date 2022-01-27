// Package atlas implements functions to work with MUD world map.
//
// Better name for this package would be "map" but it clashes with data structure name.
package atlas

import (
	"go.uber.org/zap"
	"mudbot/botutil"
)

type Atlas struct {
	lastRoom Room

	Rooms            map[int64]Room
	Coordinates      Coordinates
	nextRoomId       int64
	roomsByShorthand map[string][]Room

	movements []Direction

	server *Server

	logger *zap.SugaredLogger
}

func NewAtlas(startServer bool) *Atlas {
	a := Atlas{
		Rooms:            make(map[int64]Room),
		nextRoomId:       1,
		roomsByShorthand: make(map[string][]Room),
		logger:           botutil.NewLogger("atlas"),
	}

	a.server = NewServer(func() (map[int64]Room, Coordinates) {
		return a.Rooms, a.Coordinates
	})

	if startServer {
		a.server.Start(a)
	}

	return &a
}

func (a *Atlas) RecordMovement(dir Direction) {
	a.movements = append(a.movements, dir)
	a.logger.Debugf("Recorded movement: %v", dir)
}

func (a *Atlas) RecordRoom(room Room) {
	var hasFrom bool
	var from Direction
	fromStr := "not moving"
	if len(a.movements) > 0 {
		from = a.movements[0].Opposite()
		a.Coordinates.AddDir(a.movements[0])
		fromStr = from.String()
		a.movements = a.movements[1:]
		hasFrom = true
	}
	room.Coordinates = a.Coordinates
	var sh string
	if !room.PartialInfo {
		sh = room.Shorthand()
	}
	a.logger.Debugf("Recorded room (partial info: %t, moved from %v, shorthand %v)\n%+v", room.PartialInfo, fromStr, sh, room.Name)

	if !room.PartialInfo {
		realRoomsByShorthand, hasRealRoomsByShorthand := a.roomsByShorthand[room.Shorthand()]
		createRoom := !hasRealRoomsByShorthand

		if hasRealRoomsByShorthand {
			if len(realRoomsByShorthand) == 1 {
				a.logger.Debugf("Found single room fitting shorthand")
				foundRoomIdEqualsLastRoomId := a.lastRoom.Id == realRoomsByShorthand[0].Id
				if hasFrom && foundRoomIdEqualsLastRoomId {
					// New room ID equals previous room ID
					// Must be duplicate rooms in zone
					a.logger.Debugf("But will create new room nevertheless")
					createRoom = true
				} else {
					room = realRoomsByShorthand[0]
				}
			} else {
				a.logger.Debugf("Found multiple Rooms fitting shorthand")
				if hasFrom && a.lastRoom.Id > 0 {
					for _, realRoom := range realRoomsByShorthand {
						if realRoom.Exits[from.Opposite()] == a.lastRoom.Id {
							room = realRoom
							break
						}
					}
					// TODO: Try to find a room by x/y/z
					if room.Id == 0 {
						createRoom = true
					}
				} // else: we have multiple rooms fitting current one and no idea where we came from — leave it
			}
		}

		if createRoom {
			room.Id = a.nextRoomId
			a.nextRoomId++
			a.Rooms[room.Id] = room
			currentRoomsByShorthand, _ := a.roomsByShorthand[room.Shorthand()]
			a.roomsByShorthand[room.Shorthand()] = append(currentRoomsByShorthand, room)

			a.logger.Debugf("Created room %+v with id %v", room.Name, room.Id)

			if hasFrom && a.lastRoom.Id > 0 {
				a.lastRoom.Exits[from.Opposite()] = room.Id
				room.Exits[from] = a.lastRoom.Id
				a.logger.Debugf("Linked room %+v with %v (%v)", a.lastRoom.Name, room.Name, from.Opposite())
			}
		}
	}

	a.lastRoom = room

	a.server.sendUpdates()
}

func (a *Atlas) RecordCannotMoveFeedback() {
	if len(a.movements) > 0 {
		a.movements = a.movements[1:]
	}
	a.logger.Debugf("Recorded cannot move feedback")
}