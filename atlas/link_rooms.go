package atlas

import "mudbot/atlas/server"

func (a *Atlas) LinkRooms(fromRoom *Room, fromRoomExit Direction, toRoom *Room, toRoomExit Direction) {
	fromRoom.Exits[fromRoomExit] = toRoom.Id
	toRoom.Exits[toRoomExit] = fromRoom.Id
	a.server.SendData("link_rooms")
}

func (a *Atlas) onLinkRooms(cmd server.LinkRoomsCommand) {
	fromRoom, hasFromRoom := a.Rooms[int64(cmd.FromRoomId)]
	if !hasFromRoom {
		a.logger.Debugf("No from room with ID %v", cmd.FromRoomId)
		return
	}

	toRoom, hasToRoom := a.Rooms[int64(cmd.ToRoomId)]
	if !hasToRoom {
		a.logger.Debugf("No room with ID %v", cmd.ToRoomId)
		return
	}

	fromRoomDirection, okDirectionFrom := NewDirection(cmd.FromRoomExit)
	if !okDirectionFrom {
		a.logger.Errorf("Unknown FromRoomExit: %v", cmd.FromRoomExit)
		return
	}

	toRoomDirection, okDirectionTo := NewDirection(cmd.ToRoomExit)
	if !okDirectionTo {
		a.logger.Errorf("Unknown ToRoomExit: %v", cmd.ToRoomExit)
		return
	}

	a.LinkRooms(fromRoom, fromRoomDirection, toRoom, toRoomDirection)
}
