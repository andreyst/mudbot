package atlas

// While waiting for generics in stable...
func getFirstRoom(rooms Rooms) (firstRoom *Room) {
	for _, firstRoom = range rooms {
		break
	}

	return
}
