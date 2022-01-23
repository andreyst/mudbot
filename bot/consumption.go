package bot

import (
	"mudbot/botutil"
)

var prHungry = "Вы хотите есть."
var prThirsty = "Вы хотите пить."
var prHungryAlt = "Вам пора подкрепиться."
var prThirstyAlt = "Вы ощущаете сухость во рту."
var prHungryAndThirsty = "Вы хотите есть. Вы хотите пить."
var prHungryAndThirstyAlt = "Вам пора подкрепиться. Вы ощущаете сухость во рту."
var prTimeToEat = "Пришла пора подкрепиться."
var prTimeToDrink = "Во рту у Вас пересохло."
var prSatiated = "Вы наелись."
var prQuenched = "Вы больше не чувствуете жажды."
var prFull = "Ваш желудок не выдержит столько еды."
var prOverflowing = "Ваш желудок не выдержит столько жидкости."
var prAte = "Вы съели"
var prDrank = "Вы выпили"

func (b *Bot) ParseConsumption(s string) (res []Event) {
	if botutil.HasAnyLinePrefix(s, []string{prHungry, prHungryAlt, prHungryAndThirsty, prHungryAndThirstyAlt, prTimeToEat}) {
		b.Char.IsHungry = true
	}
	if botutil.HasAnyLinePrefix(s, []string{prThirsty, prThirstyAlt, prHungryAndThirsty, prHungryAndThirstyAlt, prTimeToDrink}) {
		b.Char.IsThirsty = true
	}
	if botutil.HasAnyLinePrefix(s, []string{prSatiated, prFull}) {
		b.Char.IsHungry = false
	}
	if botutil.HasAnyLinePrefix(s, []string{prQuenched, prOverflowing}) {
		b.Char.IsThirsty = false
	}
	if botutil.HasAnyLinePrefix(s, []string{prAte, prFull}) {
		res = append(res, EVENT_ATE)
	}
	if botutil.HasAnyLinePrefix(s, []string{prDrank, prOverflowing}) {
		res = append(res, EVENT_DRANK)
	}

	return res
}
