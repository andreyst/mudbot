package atlas

import "mudbot/atlas/server"

func (a *Atlas) DeleteRoom(roomId int64) {
	room, hasRoom := a.Rooms[roomId]
	if !hasRoom {
		return
	}

	if a.lastRoom != nil && a.lastRoom.Id == room.Id {
		a.lastRoom = nil
	}

	delete(a.roomsByShorthand[room.Shorthand()], roomId)
	delete(a.roomsByCoordinates, room.Coordinates)
	a.deleteExitFromAllRooms(room.Id)
	delete(a.Rooms, roomId)
	a.server.SendData()
}

func (a *Atlas) deleteExitFromAllRooms(roomId int64) {
	for _, room := range a.Rooms {
		for dir, exitRoomId := range room.Exits {
			if exitRoomId == roomId {
				room.Exits[dir] = 0
			}
		}
	}
}

func (a *Atlas) onDeleteRoom(cmd server.DeleteRoomCommand) {
	a.DeleteRoom(int64(cmd.RoomId))
}
