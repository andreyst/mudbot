// Package atlas implements functions to work with MUD world map.
//
// Better name for this package would be "map" but it clashes with data structure name.
package atlas

import (
	"go.uber.org/zap"
	"mudbot/botutil"
)

type Atlas struct {
	lastRoom *Room

	Rooms              map[int64]*Room
	Coordinates        Coordinates
	nextRoomId         int64
	roomsByShorthand   map[string][]*Room
	roomsByCoordinates map[Coordinates]*Room

	movements []Direction

	server *Server

	logger *zap.SugaredLogger
}

func NewAtlas(startServer bool) *Atlas {
	a := Atlas{
		Rooms:              make(map[int64]*Room),
		nextRoomId:         1,
		roomsByShorthand:   make(map[string][]*Room),
		roomsByCoordinates: make(map[Coordinates]*Room),
		logger:             botutil.NewLogger("atlas"),
	}

	a.server = NewServer(func() (map[int64]*Room, Coordinates, *Room) {
		return a.Rooms, a.Coordinates, a.lastRoom
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

func (a *Atlas) RecordRoom(room *Room) {
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
		var createRoom bool
		var lookByShorthand bool

		if a.lastRoom != nil && !a.lastRoom.IsAstraying {
			roomByCoordinates, hasRoomByCoordinates := a.roomsByCoordinates[a.Coordinates]

			if !hasRoomByCoordinates {
				roomByShorthand, hasRoomByShorthand := a.findByShorthand(room.Shorthand())
				if hasRoomByShorthand && roomByShorthand.IsAstraying {
					a.logger.Errorf("No room by coordinates but found astraying room by shorthand")
					room = roomByShorthand
					a.Coordinates = room.Coordinates
				}
				createRoom = true
			} else if room.Shorthand() != roomByCoordinates.Shorthand() {
				a.logger.Errorf("Room by coordinates does not match current room by shorthand")
				lookByShorthand = true
			} else {
				room = roomByCoordinates
			}
		} else {
			lookByShorthand = true
		}

		if lookByShorthand {
			roomByShorthand, hasRoomByShorthand := a.findByShorthand(room.Shorthand())
			if !hasRoomByShorthand {
				a.logger.Debugf("Did not find room by shorthand")
				createRoom = true
			} else {
				room = roomByShorthand
				a.Coordinates = room.Coordinates
				if a.lastRoom != nil {
					if !room.IsAstraying {
						a.logger.Debugf("Marking current room as astraying because moved from previous room not by coords")
						room.IsAstraying = true
					}
					a.lastRoom.Exits[from.Opposite()] = room.Id
				}
			}
		}

		if createRoom {
			room.Id = a.nextRoomId
			a.nextRoomId++
			a.Rooms[room.Id] = room
			currentRoomsByShorthand, _ := a.roomsByShorthand[room.Shorthand()]
			a.roomsByShorthand[room.Shorthand()] = append(currentRoomsByShorthand, room)
			a.roomsByCoordinates[room.Coordinates] = room

			a.logger.Debugf("Created room %+v with id %v", room.Name, room.Id)

			if hasFrom && a.lastRoom != nil {
				a.lastRoom.Exits[from.Opposite()] = room.Id
				if _, ok := room.Exits[from]; ok {
					room.Exits[from] = a.lastRoom.Id
				} else {
					a.logger.Debugf("Marking new and previous rooms %v as astraying because of mismatching exits", room.Id)
					room.IsAstraying = true
					a.lastRoom.IsAstraying = true
					a.logger.Debugf("%+v", a.Rooms[room.Id])
				}
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

func (a *Atlas) findByShorthand(shorthand string) (res *Room, hasRoomByShorthand bool) {
	roomsByShorthand, hasRoomByShorthand := a.roomsByShorthand[shorthand]
	if hasRoomByShorthand {
		res = roomsByShorthand[0]
		if len(roomsByShorthand) > 1 {
			a.logger.Errorf("Multiple rooms by shorthand found!")
		}
	}

	return
}
