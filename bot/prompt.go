package bot

import (
	"mudbot/botutil"
	"strconv"
	"strings"

	"github.com/oriser/regroup"
)

var promptRe = regroup.MustCompile(`(?ms)(?P<Health>-?\d+)ж (?P<Stamina>-?\d+)б (?P<ExperienceTNL>\d+)о (?P<Money>\d+)м .*Выходы:(?P<Exits>\pL*)`)

func (b *Bot) ParsePrompt(s string) Event {
	lastCr := strings.LastIndex(s, "\n")
	if lastCr >= 0 {
		s = s[lastCr+1:]
	}
	s = botutil.StripAnsi(s)

	promptMatch, _ := promptRe.Groups(s)
	if promptMatch == nil {
		return EVENT_NOP
	}

	// TODO: Add proper error handling
	b.Char.Health, _ = strconv.Atoi(promptMatch["Health"])
	b.Char.Stamina, _ = strconv.Atoi(promptMatch["Stamina"])
	b.Char.ExperienceTNL, _ = strconv.Atoi(promptMatch["ExperienceTNL"])
	b.Char.Money, _ = strconv.Atoi(promptMatch["Money"])

	b.ParseFight(s)

	return EVENT_PROMPT
}
