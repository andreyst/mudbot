//go:generate stringer -type=Event -output events_string.go

package mud

type Event int

const (
	EVENT_TICK Event = iota
	EVENT_NOP

	EVENT_ENCODING_MENU
	EVENT_LOGIN_PROMPT
	EVENT_PASSWORD_PROMPT
	EVENT_PRESS_CRFL_PROMPT
	EVENT_GAME_MENU

	EVENT_PROMPT
	EVENT_ROOM

	EVENT_SCORE

	EVENT_ATE
	EVENT_DRANK

	EVENT_RESTING
	EVENT_SAT
	EVENT_STOOD_UP
	EVENT_WENT_TO_SLEEP
	EVENT_WOKE_UP

	EVENT_CANNOT_MOVE_IN_THIS_DIRECTION
	EVENT_CANNOT_BECAUSE_RESTING
	EVENT_NO_SUCH_ITEM
	EVENT_WATER_CONTAINER_EMPTY
)
