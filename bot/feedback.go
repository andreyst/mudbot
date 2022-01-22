package bot

import (
	"mudbot/botutil"
	"regexp"
)

var prNoSuchItem = "Здесь нет такого предмета."
var prNoSuchItemInInventory = "У Вас в руках нет такого предмета."
var itemEmptyMatcher = regexp.MustCompile(`(?ms)^В .* пусто.$`)
var prCantBecauseSleeping = "Сделать это во сне, что ли?"
var prCantBecauseSitting = "Для этого необходимо встать на ноги."
var prCantBecauseResting = "Вы слишком расслаблены, и Вам сейчас не до этого."

func (b *Bot) ParseFeedback(s string) (res []Event) {
	if botutil.HasAnyLinePrefix(s, []string{prNoSuchItem, prNoSuchItemInInventory}) {
		res = append(res, EVENT_NO_SUCH_ITEM)
	}
	if itemEmptyMatcher.Match([]byte(s)) {
		res = append(res, EVENT_WATER_CONTAINER_EMPTY)
	}
	if botutil.HasAnyLinePrefix(s, []string{prCantBecauseResting}) {
		b.Char.Position = POSITION_RESTING
	}
	if botutil.HasAnyLinePrefix(s, []string{prCantBecauseSitting}) {
		b.Char.Position = POSITION_SITTING
	}
	if botutil.HasAnyLinePrefix(s, []string{prCantBecauseSleeping}) {
		b.Char.Position = POSITION_SLEEPING
	}

	return
}
