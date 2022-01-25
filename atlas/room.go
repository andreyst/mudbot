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

func NewRoom() Room {
	r := Room{
		Exits: make(map[Direction]int64),
	}

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
