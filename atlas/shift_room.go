package atlas

import "mudbot/atlas/server"

func (a *Atlas) ShiftRoom(roomId int64, direction Direction) {
	if _, roomExists := a.Rooms[roomId]; !roomExists {
		a.logger.Errorf("Room %v does not exist", roomId)
		return
	}

	a.logger.Debugf("Shifting room %v %v", roomId, direction)
	visited := make(map[int64]bool)
	a.doShift(roomId, direction, &visited)
	a.server.SendData("shift_room")
}

func (a *Atlas) doShift(roomId int64, direction Direction, visited *map[int64]bool) {
	if (*visited)[roomId] {
		return
	}
	(*visited)[roomId] = true

	room := a.Rooms[roomId]
	for exitDir, exitRoomId := range room.Exits {
		if exitDir == direction.Opposite() || exitRoomId == 0 {
			continue
		}
		// TODO: Fix shifting when side links are not in the same axis
		if exitDir != direction || room.Distance(*a.Rooms[exitRoomId]) == 1 {
			// Do not move in required direction if there is space
			a.doShift(exitRoomId, direction, visited)
		}
	}

	delete(a.roomsByCoordinates, room.Coordinates)

	newCoordinates := room.Coordinates
	newCoordinates.Shift(direction)
	if roomAtNewCoordinates, hasRoomAtNewCoordinates := a.roomsByCoordinates[newCoordinates]; hasRoomAtNewCoordinates {
		a.doShift(roomAtNewCoordinates.Id, direction, visited)
	}

	room.Shift(direction)
	a.roomsByCoordinates[room.Coordinates] = room
}

func (a *Atlas) onShiftRoom(cmd server.ShiftRoomCommand) {
	dir, ok := NewDirection(cmd.Direction)
	if !ok {
		a.logger.Infof("Wrong direction in shift command, expected one of directions in NewDirection(): %v", cmd.Direction)
		return
	}

	a.ShiftRoom(int64(cmd.RoomId), dir)
}
