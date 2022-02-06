package atlas

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"mudbot/botutil"
	"sort"
	"strings"
)

type Exits map[Direction]int64

type Room struct {
	Id          int64
	Coordinates Coordinates

	Name        string
	Description string
	PartialInfo bool

	IsTricky bool

	Exits Exits
	Items []string
	Mobs  []string
}

func NewEmptyRoom() *Room {
	r := Room{
		Exits: make(Exits),
	}

	return &r
}

func NewRoom(name string) *Room {
	r := NewEmptyRoom()
	r.Name = name

	return r
}

func NewRoomWithExits(name string, exits []Direction) *Room {
	r := NewRoom(name)
	for _, exit := range exits {
		r.Exits[exit] = 0
	}
	return r
}

func NewFilledRoom(name string, description string, exits []Direction, items []string, mobs []string) *Room {
	r := NewRoomWithExits(name, exits)

	r.Description = description
	r.Items = items
	r.Mobs = mobs

	return r
}

func (r Room) Shorthand() string {
	if r.PartialInfo {
		panic(errors.New("cannot calculate shorthand for room with partial info"))
	}

	exitsStr := make([]string, len(r.Exits))
	for dir := range r.Exits {
		exitsStr = append(exitsStr, dir.String())
	}
	sort.Strings(exitsStr)

	hash := sha1.Sum([]byte(r.Name + "_" + r.Description + "_" + strings.Join(exitsStr, "_")))
	return hex.EncodeToString(hash[:])
}

func (r Room) Distance(to Room) int64 {
	distanceX := botutil.Abs(r.Coordinates.X - to.Coordinates.X)
	distanceY := botutil.Abs(r.Coordinates.Y - to.Coordinates.Y)
	distanceZ := botutil.Abs(r.Coordinates.Z - to.Coordinates.Z)
	return distanceX + distanceY + distanceZ
}

func (r *Room) Shift(direction Direction) {
	r.Coordinates.Shift(direction)
}
