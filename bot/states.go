//go:generate stringer -type=State -output states_string.go

package bot

type State int

const (
	STATE_INITIALIZING State = iota

	STATE_IDLE
	STATE_STUCK

	STATE_DRINKING
	STATE_EATING

	STATE_RESTING
)
