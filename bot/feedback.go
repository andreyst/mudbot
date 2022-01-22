package bot

import (
	"mudbot/botutil"
	"regexp"
)

var prNoSuchItem = "Здесь нет такого предмета."
var prNoSuchItemInInventory = "У Вас в руках нет такого предмета."
var itemEmptyMatcher = regexp.MustCompile(`(?ms)^В .* пусто.$`)

func (b *Bot) ParseFeedback(s string) (res []Event) {
	if botutil.HasAnyLinePrefix(s, []string{prNoSuchItem, prNoSuchItemInInventory}) {
		res = append(res, EVENT_NO_SUCH_ITEM)
	}
	if itemEmptyMatcher.Match([]byte(s)) {
		res = append(res, EVENT_WATER_CONTAINER_EMPTY)
	}

	return
}
