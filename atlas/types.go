//go:generate stringer -type=Direction -output types_string.go

package atlas

type Direction int

const (
	DIRECTION_NORTH Direction = iota
	DIRECTION_SOUTH
	DIRECTION_WEST
	DIRECTION_EAST
	DIRECTION_UP
	DIRECTION_DOWN
	DIRECTION_ENTER
	DIRECTION_EXIT
)

type Room struct {
	Name        string
	Description string
	NoInfo      bool
	Exits       []string
	Items       []string
	Mobs        []string
}
