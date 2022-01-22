//go:generate stringer -type=Event -output events_string.go

package bot

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

	EVENT_SCORE

	EVENT_ATE
	EVENT_DRANK

	EVENT_NO_SUCH_ITEM
	EVENT_WATER_CONTAINER_EMPTY
)
