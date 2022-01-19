package bot

import (
	"mudbot/botutil"
	"strconv"
	"strings"

	"github.com/oriser/regroup"
)

var promptRe = regroup.MustCompile(`(?ms)(?P<Health>-?\d+)ж (?P<Stamina>-?\d+)б (?P<ExperienceTNL>\d+)о (?P<Money>\d+)м .*Выходы:(?P<Exits>\pL*)`)

func (b *Bot) ParsePrompt(s string) bool {
	lastCr := strings.LastIndex(s, "\n")
	if lastCr >= 0 {
		s = s[lastCr+1:]
	}
	s = botutil.StripAnsi(s)

	promptMatch, _ := promptRe.Groups(s)
	if promptMatch == nil {
		return false
	}

	// TODO: Add proper error handling
	b.char.Health, _ = strconv.Atoi(promptMatch["Health"])
	b.char.Stamina, _ = strconv.Atoi(promptMatch["Stamina"])
	b.char.ExperienceTNL, _ = strconv.Atoi(promptMatch["ExperienceTNL"])
	b.char.Money, _ = strconv.Atoi(promptMatch["Money"])

	b.ParseFight(s)

	b.logger.Debugf("char after ParsePrompt: %+v", b.char)
	b.logger.Debugf("fight after ParsePrompt: %+v", b.fight)

	return false
}
