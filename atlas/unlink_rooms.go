package atlas

import "mudbot/atlas/server"

func (a *Atlas) UnlinkRooms(fromRoom, toRoom *Room) {
	removeExitsToRoom(fromRoom, toRoom.Id)
	removeExitsToRoom(toRoom, fromRoom.Id)
	a.server.SendData("unlink_rooms")
}

func removeExitsToRoom(room *Room, toRoomId int64) {
	for dir, roomId := range room.Exits {
		if roomId == toRoomId {
			room.Exits[dir] = 0
		}
	}
}

func (a *Atlas) onUnlinkRooms(cmd server.UnlinkRoomsCommand) {
	fromRoom, hasFromRoom := a.Rooms[int64(cmd.FromRoomId)]
	if !hasFromRoom {
		a.logger.Debugf("No from room with ID %v", fromRoom.Id)
		return
	}

	toRoom, hasToRoom := a.Rooms[int64(cmd.ToRoomId)]
	if !hasToRoom {
		a.logger.Debugf("No room with ID %v", toRoom.Id)
		return
	}

	a.UnlinkRooms(fromRoom, toRoom)
}
