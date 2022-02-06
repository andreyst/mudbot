package atlas

func (a *Atlas) dataProvider() interface{} {
	data := struct {
		Rooms       Rooms
		Coordinates Coordinates
		Room        *Room
	}{
		a.Rooms,
		a.Coordinates,
		a.lastRoom,
	}

	return data
}
