package atlas

func (a *Atlas) dataProvider() interface{} {
	data := struct {
		Rooms       map[int64]*Room
		Coordinates Coordinates
		Room        *Room
	}{
		a.Rooms,
		a.Coordinates,
		a.lastRoom,
	}

	return data
}
