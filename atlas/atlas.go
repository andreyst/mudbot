// Package atlas implements functions to work with MUD world map.
//
// Better name for this package would be "map" but it clashes with data structure name.
package atlas

import (
	"go.uber.org/zap"
	"mudbot/atlas/server"
	"mudbot/botutil"
)

type Atlas struct {
	lastRoom *Room

	Rooms              map[int64]*Room
	Coordinates        Coordinates
	nextRoomId         int64
	roomsByShorthand   map[string]map[int64]*Room
	roomsByCoordinates map[Coordinates]*Room

	movements []Direction

	server *server.Server

	logger *zap.SugaredLogger
}

func NewAtlas() *Atlas {
	a := Atlas{
		Rooms:              make(map[int64]*Room),
		nextRoomId:         1,
		roomsByShorthand:   make(map[string]map[int64]*Room),
		roomsByCoordinates: make(map[Coordinates]*Room),
		logger:             botutil.NewLogger("atlas"),
	}

	a.server = server.NewServer(a.dataProvider)
	a.server.OnShiftRoom = a.onShiftRoom
	a.server.OnDeleteRoom = a.onDeleteRoom

	return &a
}

func (a *Atlas) StartServer() {
	a.server.Start()
}

func (a *Atlas) RecordMovement(dir Direction) {
	a.movements = append(a.movements, dir)
	a.logger.Debugf("Recorded movement: %v", dir)
}

func (a *Atlas) RecordRelocation() {
	a.lastRoom = nil
	a.logger.Debugf("Recorded relocation")
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

		if a.lastRoom == nil {
			// There was no last room - if there are no rooms at all,
			// lets create first one. In other case, let's try to locate ourselves,
			// but do not create new rooms not connected to main map.
			if len(a.Rooms) == 0 {
				a.logger.Debugf("No rooms, creating first room")
				createRoom = true
			} else {
				roomsByShorthand := a.roomsByShorthand[room.Shorthand()]
				if len(roomsByShorthand) == 1 {
					a.logger.Debugf("Found single room by shorthand, using it as current location")
					room = getFirstRoom(roomsByShorthand)
				} else {
					a.logger.Debugf("Cannot find single room by shorthand, move to unambigious and record room to locate self")
				}
			}
		} else if !hasFrom {
			// There was no movement.
			if a.lastRoom.Shorthand() == room.Shorthand() {
				// Current room is the same as last room – as expected.
				room = a.lastRoom
			} else {
				// For some reason current room does not match last room without movement
				// Let's update current location.
				a.logger.Debugf("Did not move, but current room shorthand does not equal last room - need to update")
				roomsByShorthand := a.roomsByShorthand[room.Shorthand()]
				if len(a.roomsByShorthand) == 1 {
					a.logger.Debugf("Found single room by shorthand, using it as current location")
					room = getFirstRoom(roomsByShorthand)
				} else {
					a.logger.Debugf("Cannot find single room by shorthand, move to unambigious and record room to locate self")
				}
			}
		} else {
			// Do we have a linked exit from last room corresponding to where we went?
			roomIdByExit, hasExit := a.lastRoom.Exits[from.Opposite()]
			if !hasExit || roomIdByExit == 0 {
				// There is no linked exit. Was last room tricky?
				if !a.lastRoom.IsTricky {
					// Last room was not tricky. Let's try to find by coordinates
					// or by shorthand with more strict matching rules.
					a.logger.Debugf("No linked exit, looking for candidate room")
					roomByCoordinates, hasRoomByCoordinates := a.roomsByCoordinates[a.Coordinates]
					var checkByShorthand bool
					if hasRoomByCoordinates && roomByCoordinates.Shorthand() == room.Shorthand() {
						// There is such a room by coordinates. Let's check
						// if its corresponding exit is already linked.
						a.logger.Debugf("Found room by coordinates")
						roomByCoordinatesExit := roomByCoordinates.Exits[from]
						if roomByCoordinatesExit == 0 {
							// It does not have a linked corresponding exit. Let's use it.
							room = roomByCoordinates
							a.Coordinates = room.Coordinates
							a.lastRoom.Exits[from.Opposite()] = room.Id
							exitRoomId, hasExit := room.Exits[from]
							if hasExit && exitRoomId == 0 {
								room.Exits[from] = a.lastRoom.Id
							}
						} else {
							// It already has a linked corresponding exit. Let's create a new room.
							a.logger.Debugf("But it already has linked corresponding exit, creating new")
							checkByShorthand = true
						}
					} else {
						checkByShorthand = true
					}

					if checkByShorthand {
						// Let's see if there are rooms by shorthand.
						// Due to duplicate rooms there is a high chance
						// that one of them will be picked by shorthand instead of
						// creating new. To mitigate this, let's say that
						// there should be only one room by shorthand.
						roomsByShorthand, hasRoomByShorthand := a.roomsByShorthand[room.Shorthand()]
						if hasRoomByShorthand && len(roomsByShorthand) == 1 {
							roomByShorthand := getFirstRoom(roomsByShorthand)
							_, hasRoomByShorthandExit := roomByShorthand.Exits[from]
							if hasRoomByShorthandExit ||
								roomByShorthand.Id == a.lastRoom.Id {
								a.logger.Debugf("No fitting room by shorthand or coordinates found, creating new")
								createRoom = true
							} else {
								a.logger.Debugf("Found room by shorthand, using it")
								room = roomByShorthand
								a.Coordinates = room.Coordinates
								a.lastRoom.Exits[from.Opposite()] = room.Id
							}
						} else {
							a.logger.Debugf("No fitting room by shorthand or coordinates found, creating new")
							createRoom = true
						}
					}
				} else {
					// Last room was tricky. Look for new room by shorthand
					// with more relaxed matching rules - because last room was tricky,
					// it could lead anywhere.
					roomsByShorthand, hasRoomByShorthand := a.roomsByShorthand[room.Shorthand()]
					if hasRoomByShorthand && len(roomsByShorthand) == 1 {
						a.logger.Debugf("Found room by shorthand, using it")
						room = getFirstRoom(roomsByShorthand)
						a.Coordinates = room.Coordinates
						a.lastRoom.Exits[from.Opposite()] = room.Id

						if !room.IsTricky {
							a.logger.Debugf("Marking it as tricky because moved from tricky room not by coords")
							room.IsTricky = true
						}
					} else {
						a.logger.Debugf("No room by shorthand found, creating new")
						createRoom = true
					}
				}
			} else {
				// There is a linked exit from last room. Does the room linked to it has the same
				// shorthand as the one we moved to?
				roomByExit := a.Rooms[roomIdByExit]
				if roomByExit.Shorthand() == room.Shorthand() {
					// Yes, we went to the same room as we went to. Let's use it.
					a.logger.Debugf("Found room by exit, using it")
					room = roomByExit
					a.Coordinates = room.Coordinates
				} else {
					// No, rooms are different. Could be a random exit,
					// an exit that leads back not to where you came from,
					// exit changed due to trigger, relocation event, and so on.
					// Let's try to find if we already know about the room we went to.
					a.logger.Errorf("Room by exit does not equal actual room")

					// TODO: Redo for multiple rooms by shorthand
					roomsByShorthand, hasRoomsByShorthand := a.roomsByShorthand[room.Shorthand()]
					if !hasRoomsByShorthand {
						// There is no such room. Let's create it.
						a.logger.Debugf("Did not find room by shorthand")
						createRoom = true
					} else if len(roomsByShorthand) == 1 {
						a.logger.Debugf("Found room by shorthand, marking last room as tricky because existing exit lead to unexpected room")
						room = getFirstRoom(roomsByShorthand)
						a.Coordinates = room.Coordinates
						a.lastRoom.Exits[from.Opposite()] = room.Id
						roomExit, roomHasExit := room.Exits[from.Opposite()]
						if roomHasExit && roomExit > 0 {
							room.Exits[from.Opposite()] = a.lastRoom.Id
						}
					} else {
						a.logger.Debugf("Room by shorthand was not tricky")
						createRoom = true
					}

					// There is such room. Does it have an unlinked exit corresponding
					// with movement direction
					//roomByShorthandExit, hasRoomByShorthandExit := roomByShorthand.Exits[from]
					//if !hasRoomByShorthandExit {
					//	// There is no exit at all. Is this room we went to tricky?
					//} else if roomByShorthandExit == 0 {
					//	// There is no linked corresponding exit - lets think that we moved to this room.
					//	// TODO: Redo to picking closest room by shorthand
					//	//var closestRoomByShorthand *Room
					//	//for _, roomByShorthand := range roomsByShorthand {
					//	//	roomByShorthandExit := roomByShorthand.Exits[from]
					//	//	if roomByShorthandExit > 0 {
					//	//		continue
					//	//	}
					//	//	if closestRoomByShorthand == nil ||
					//	//		a.lastRoom.Distance(*closestRoomByShorthand) > a.lastRoom.Distance(*roomByShorthand) {
					//	//		closestRoomByShorthand = roomByShorthand
					//	//	}
					//	//}
					//
					//	a.logger.Debugf("Found room by shorthand without linked corresponding exit")
					//	room = roomByShorthand
					//	a.Coordinates = room.Coordinates
					//	a.lastRoom.Exits[from.Opposite()] = room.Id
					//	// Not linking back, though.
					//} else {
					//	// There is a linked exit.
					//	// TODO: What about multiple rooms leading to single non-adjacent room? Or via one ways.
					//	a.logger.Debugf("Room by shorthand did not have empty corresponding exit")
					//	createRoom = true
					//}
					//}
				}
			}
		}

		// TODO: Remove previous version
		//if a.lastRoom != nil && !a.lastRoom.IsTricky {
		//	roomByCoordinates, hasRoomByCoordinates := a.roomsByCoordinates[a.Coordinates]
		//
		//	if !hasRoomByCoordinates {
		//		roomByShorthand, hasRoomByShorthand := a.findByShorthand(room.Shorthand())
		//		if hasRoomByShorthand && roomByShorthand.IsTricky {
		//			a.logger.Errorf("No room by coordinates but found tricky room by shorthand")
		//			room = roomByShorthand
		//			a.Coordinates = room.Coordinates
		//		}
		//		createRoom = true
		//	} else if room.Shorthand() != roomByCoordinates.Shorthand() {
		//		a.logger.Errorf("Room by coordinates does not match current room by shorthand")
		//		lookByShorthand = true
		//	} else {
		//		room = roomByCoordinates
		//	}
		//} else {
		//	lookByShorthand = true
		//}
		//
		//if lookByShorthand {
		//	roomByShorthand, hasRoomByShorthand := a.findByShorthand(room.Shorthand())
		//	if !hasRoomByShorthand {
		//		a.logger.Debugf("Did not find room by shorthand")
		//		createRoom = true
		//	} else {
		//		room = roomByShorthand
		//		a.Coordinates = room.Coordinates
		//		if a.lastRoom != nil {
		//			if !room.IsTricky {
		//				a.logger.Debugf("Marking current room as tricky because moved from previous room not by coords")
		//				room.IsTricky = true
		//			}
		//			a.lastRoom.Exits[from.Opposite()] = room.Id
		//		}
		//	}
		//}

		if createRoom {
			existingRoomByCoordinates, hasExistingRoomByCoordinates := a.roomsByCoordinates[room.Coordinates]
			if hasExistingRoomByCoordinates && !hasFrom {
				a.logger.Errorf("Not creating room becuase it clashes with existing by coordinates, but have no movement to figure out where to shift it")
			} else {
				if hasExistingRoomByCoordinates {
					a.ShiftRoom(existingRoomByCoordinates.Id, from.Opposite())
				}

				room.Id = a.nextRoomId
				a.nextRoomId++
				a.Rooms[room.Id] = room
				if a.roomsByShorthand[room.Shorthand()] == nil {
					a.roomsByShorthand[room.Shorthand()] = make(map[int64]*Room)
				}
				a.roomsByShorthand[room.Shorthand()][room.Id] = room
				a.roomsByCoordinates[room.Coordinates] = room

				a.logger.Debugf("Created room %+v with id %v", room.Name, room.Id)

				if hasFrom && a.lastRoom != nil {
					a.lastRoom.Exits[from.Opposite()] = room.Id
					if _, ok := room.Exits[from]; ok {
						room.Exits[from] = a.lastRoom.Id
					} else {
						a.logger.Debugf("Marking new and previous rooms %v as tricky because of mismatching exits", room.Id)
						room.IsTricky = true
						a.lastRoom.IsTricky = true
						a.logger.Debugf("%+v", a.Rooms[room.Id])
					}
					a.logger.Debugf("Linked room %+v with %v (%v)", a.lastRoom.Name, room.Name, from.Opposite())
				}
			}
		}
	}

	if room.Id > 0 {
		a.lastRoom = room
	} else {
		a.lastRoom = nil
	}

	a.server.SendData()
}

func (a *Atlas) RecordCannotMoveFeedback() {
	if len(a.movements) > 0 {
		a.movements = a.movements[1:]
	}
	a.logger.Debugf("Recorded cannot move feedback")
}
