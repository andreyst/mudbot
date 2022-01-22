package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCombatant1(t *testing.T) {
	s := `хомячок:умирает`

	ref := Combatant{
		Name:      "хомячок",
		Condition: CONDITION_DYING,
		Position:  POSITION_STANDING,
	}

	b := NewBot(Credentials{})
	c := b.ParseCombatant(s)

	assert.Equal(t, ref, c)
}

func TestParseCombatant2(t *testing.T) {
	s := `хомячок:умирает,леж`

	ref := Combatant{
		Name:      "хомячок",
		Condition: CONDITION_DYING,
		Position:  POSITION_LYING,
	}

	b := NewBot(Credentials{})
	c := b.ParseCombatant(s)

	assert.Equal(t, ref, c)
}

func TestParseCombatant3(t *testing.T) {
	s := `хомячок:великолепное,леж`

	ref := Combatant{
		Name:      "хомячок",
		Condition: CONDITION_PERFECT,
		Position:  POSITION_LYING,
	}

	b := NewBot(Credentials{})
	c := b.ParseCombatant(s)

	assert.Equal(t, ref, c)
}

func TestFight1(t *testing.T) {
	s := ``

	ref := Fight{}

	b := NewBot(Credentials{})
	b.ParseFight(s)

	assert.Equal(t, ref, b.Fight)
}

func TestFight2(t *testing.T) {
	s := `[хомячок:умирает,сид]`

	ref := Fight{
		IsActive: true,
		Enemy: Combatant{
			Name:      "хомячок",
			Condition: CONDITION_DYING,
			Position:  POSITION_SITTING,
		},
	}

	b := NewBot(Credentials{})
	b.ParseFight(s)

	assert.Equal(t, ref, b.Fight)
}

func TestFight3(t *testing.T) {
	s := `[Венге:великолепное,леж] [хомячок:среднее]`

	ref := Fight{
		IsActive: true,
		Enemy: Combatant{
			Name:      "хомячок",
			Condition: CONDITION_AVERAGE,
			Position:  POSITION_STANDING,
		},
	}

	b := NewBot(Credentials{})
	b.ParseFight(s)

	assert.Equal(t, POSITION_LYING, b.Char.Position)
	assert.Equal(t, ref, b.Fight)
}

func TestFight4(t *testing.T) {
	s := `[Венге:великолепное] [Ярра:великолепное,леж] [хомячок:среднее]`

	ref := Fight{
		IsActive: true,
		HaveTank: true,
		Tank: Combatant{
			Name:      "Ярра",
			Condition: CONDITION_PERFECT,
			Position:  POSITION_LYING,
		},
		Enemy: Combatant{
			Name:      "хомячок",
			Condition: CONDITION_AVERAGE,
			Position:  POSITION_STANDING,
		},
	}

	b := NewBot(Credentials{})
	b.ParseFight(s)

	assert.Equal(t, ref, b.Fight)
}
