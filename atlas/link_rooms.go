package atlas

import "mudbot/atlas/server"

func (a *Atlas) LinkRooms(fromRoomId int64, directionFrom Direction, toRoomId int64, directionTo Direction) {

}

func (a *Atlas) onLinkRooms(cmd server.LinkRoomsCommand) {
	directionFrom, okDirectionFrom := NewDirection(cmd.DirectionFrom)
	if !okDirectionFrom {
		a.logger.Errorf("Unknown DirectionFrom: %v", cmd.DirectionFrom)
		return
	}
	directionTo, okDirectionTo := NewDirection(cmd.DirectionTo)
	if !okDirectionTo {
		a.logger.Errorf("Unknown DirectionTo: %v", cmd.DirectionTo)
		return
	}

	a.LinkRooms(int64(cmd.FromRoomId), directionFrom, int64(cmd.ToRoomId), directionTo)
}
