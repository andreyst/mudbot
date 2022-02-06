package atlas

func (a *Atlas) AddRoom(room *Room, hasMoved bool, movedFrom Direction) *Room {
	existingRoomByCoordinates, hasExistingRoomByCoordinates := a.roomsByCoordinates[room.Coordinates]
	if hasExistingRoomByCoordinates && hasMoved {
		a.ShiftRoom(existingRoomByCoordinates.Id, movedFrom.Opposite())
	}

	room.Id = a.nextRoomId
	a.nextRoomId++
	a.Rooms[room.Id] = room
	if a.roomsByShorthand[room.Shorthand()] == nil {
		a.roomsByShorthand[room.Shorthand()] = make(Rooms)
	}
	a.roomsByShorthand[room.Shorthand()][room.Id] = room
	a.roomsByCoordinates[room.Coordinates] = room

	a.logger.Debugf("Added room %+v with id %v", room.Name, room.Id)

	if hasMoved && a.lastRoom != nil {
		a.lastRoom.Exits[movedFrom.Opposite()] = room.Id
		if _, ok := room.Exits[movedFrom]; ok {
			room.Exits[movedFrom] = a.lastRoom.Id
		} else {
			a.logger.Debugf("Marking new and previous rooms %v as tricky because of mismatching exits", room.Id)
			room.IsTricky = true
			a.lastRoom.IsTricky = true
			a.logger.Debugf("%+v", a.Rooms[room.Id])
		}
		a.logger.Debugf("Linked room %+v with %v (%v)", a.lastRoom.Name, room.Name, movedFrom.Opposite())
	}

	return room
}

func (a *Atlas) AddRoomWithoutMovement(room *Room) *Room {
	return a.AddRoom(room, false, DIRECTION_UNKNOWN)
}
