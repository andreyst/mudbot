package bot

import (
	"mudbot/botutil"
)

var encodingStr = "Please choose a MUD client"
var nameStr = "Введите имя Вашего персонажа или \"новый\" для создания нового:"
var passwordStr = "Пароль:"
var crLfStr = "*** НАЖМИТЕ ВВОД:"
var enterGameStr = "1) Войти в игру."
var welcomeStr = "Добро пожаловать на Кринн! Пусть Ваш визит будет увлекательным!"
var reconnectingStr = "Воссоединяемся."

func (b *Bot) ParseLogin(s string) {
	if botutil.HasLinePrefix(s, encodingStr) {
		b.SendToMud("u")
	} else if botutil.HasLinePrefix(s, nameStr) {
		b.SendToMud(b.Credentials.Login)
	} else if botutil.HasLinePrefix(s, passwordStr) {
		b.SendToMudWithoutEcho(b.Credentials.Password)
	} else if botutil.HasLinePrefix(s, crLfStr) {
		b.SendToMud("")
	} else if botutil.HasLinePrefix(s, enterGameStr) {
		b.SendToMud("1")
	} else if botutil.HasAnyLinePrefix(s, []string{welcomeStr, reconnectingStr}) {
		b.InGame = true
		b.Step()
	}
}
