package atlas

// While waiting for generics in stable...
func getFirstRoom(rooms map[int64]*Room) (firstRoom *Room) {
	for _, firstRoom = range rooms {
		break
	}

	return
}
