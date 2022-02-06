package atlas

import "mudbot/atlas/server"

func (a *Atlas) LinkRoom(fromRoom *Room, fromRoomExit Direction, toRoom *Room) {
	fromRoom.Exits[fromRoomExit] = toRoom.Id
	a.server.SendData("link_room")
}

func (a *Atlas) onLinkRoom(cmd server.LinkRoomCommand) {
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

	a.LinkRoom(fromRoom, fromRoomDirection, toRoom)
}
