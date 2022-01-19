package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrompt(t *testing.T) {
	s := "36ж 69б 3000о 30м [хомячок:умирает] Выходы:СВЮЗ>"

	refChar := Char{
		Health:        36,
		Stamina:       69,
		ExperienceTNL: 3000,
		Money:         30,
	}

	refFight := Fight{
		IsActive: true,
		Enemy: Combatant{
			Name:      "хомячок",
			Condition: CONDITION_DYING,
			Position:  POSITION_STANDING,
		},
	}

	b := NewBot()
	b.ParsePrompt(s)

	assert.Equal(t, refChar, b.char)
	assert.Equal(t, refFight, b.fight)
}

func TestParsePromptIgnoresAnsi(t *testing.T) {
	s := "\x1b[0;32m36ж\x1b[0;0m \x1b[0;32m100б\x1b[0;0m 2952о 30м Выходы:Ю>"

	b := NewBot()
	b.ParsePrompt(s)

	assert.Equal(t, 36, b.char.Health)
}

func TestParsePromptParsesOnlyLastLine(t *testing.T) {
	s := "[хомячок:умирает]\n36ж 69б 3000о 30м Выходы:СВЮЗ>"

	refFight := Fight{}

	b := NewBot()
	b.ParsePrompt(s)

	assert.Equal(t, 36, b.char.Health)
	assert.Equal(t, refFight, b.fight)
}
