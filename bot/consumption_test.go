package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConsumption(t *testing.T) {
	s := `Вы мужчина Венге, наемник 1 уровня. Вам 19 лет.
У Вас сейчас 36 из 40 очков жизни и 98 из 100 очков бодрости.
Ваша душа чиста, и Ваша вера в идеалы Добра непоколебима.
Вы привержены идеалам Добра.
Вы набрали 333 очков опыта. До следующего уровня осталось 3000.
У Вас есть 48 стальных монет.
Вы играете 20 дней и 23 часа.
Вы отдыхаете.
Вы хотите есть.`

	b := NewBot(Credentials{})
	b.ParseConsumption(s)

	assert.True(t, b.Char.IsHungry)
}
