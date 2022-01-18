package char

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("LOG_LEVEL", "DEBUG")
}
func TestParseScore1(t *testing.T) {
	s := `Вы мужчина Венге, наемник 1 уровня. Вам 19 лет.
У Вас сейчас 36 из 40 очков жизни и 98 из 100 очков бодрости.
Ваша душа чиста, и Ваша вера в идеалы Добра непоколебима.
Вы привержены идеалам Добра.
Вы набрали 333 очков опыта. До следующего уровня осталось 3000.
У Вас есть 48 стальных монет.
Вы играете 20 дней и 23 часа.
Вы отдыхаете.
Вы хотите есть.`

	ref := Stats{
		Name:             "Венге",
		Age:              19,
		Level:            1,
		Experience:       333,
		ExperienceTNL:    3000,
		MoralAlignment:   MORAL_ALIGNMENT_GOOD,
		EthicalAlignment: ETHICAL_ALIGNMENT_GOOD,
		Health:           36,
		MaxHealth:        40,
		Stamina:          98,
		MaxStamina:       100,
		Money:            48,
		IsHungry:         true,
		IsThirsty:        false,
		Position:         POSITION_RESTING,
		Gender:           GENDER_MALE,
		Class:            CLASS_MERCENARY,
	}

	c := NewChar()
	c.ParseScore(s)

	assert.Equal(t, ref, c.Stats)
}

func TestParseScoreBadInput(t *testing.T) {
	c := NewChar()
	assert.False(t, c.ParseScore("zxczxczcaqd"))
}
