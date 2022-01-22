package bot

import (
	"mudbot/botutil"
)

var encodingStr = "Please choose a MUD client"
var nameStr = "Введите имя Вашего персонажа или \"новый\" для создания нового:"
var passwordStr = "Пароль:"
var crLfStr = "*** НАЖМИТЕ ВВОД:"
var enterGameStr = "1) Войти в игру."

func (b *Bot) ParseLogin(s string) Event {
	if botutil.HasLinePrefix(s, encodingStr) {
		return EVENT_ENCODING_MENU
	} else if botutil.HasLinePrefix(s, nameStr) {
		return EVENT_LOGIN_PROMPT
	} else if botutil.HasLinePrefix(s, passwordStr) {
		return EVENT_PASSWORD_PROMPT
	} else if botutil.HasLinePrefix(s, crLfStr) {
		return EVENT_PRESS_CRFL_PROMPT
	} else if botutil.HasLinePrefix(s, enterGameStr) {
		return EVENT_GAME_MENU
	} else {
		return EVENT_NOP
	}
}
