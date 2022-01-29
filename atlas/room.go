package atlas

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

type Room struct {
	Id          int64
	Coordinates Coordinates

	Name        string
	Description string
	PartialInfo bool

	Exits map[Direction]int64
	Items []string
	Mobs  []string
}

func NewRoom() *Room {
	r := Room{
		Exits: make(map[Direction]int64),
	}

	return &r
}

func NewPrefilledRoom(name string, description string, exits []Direction, items []string, mobs []string) *Room {
	r := NewRoom()

	r.Name = name
	r.Description = description
	for _, exit := range exits {
		r.Exits[exit] = 0
	}
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
	distanceX := r.Coordinates.X - to.Coordinates.X
	distanceY := r.Coordinates.Y - to.Coordinates.Y
	distanceZ := r.Coordinates.Z - to.Coordinates.Z
	return distanceX + distanceY + distanceZ
}

func (r *Room) Shift(direction Direction) {
	r.Coordinates.Shift(direction)
}
