package atlas

import "mudbot/atlas/server"

func (a *Atlas) DeleteRoom(room *Room) {
	if a.lastRoom != nil && a.lastRoom.Id == room.Id {
		a.lastRoom = nil
	}

	delete(a.roomsByShorthand[room.Shorthand()], room.Id)
	delete(a.roomsByCoordinates, room.Coordinates)
	a.deleteExitFromAllRooms(room.Id)
	delete(a.Rooms, room.Id)
	a.server.SendData("delete_room")
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
	room, hasRoom := a.Rooms[int64(cmd.RoomId)]
	if !hasRoom {
		return
	}

	a.DeleteRoom(room)
}
