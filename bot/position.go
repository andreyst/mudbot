package bot

import (
	"mudbot/botutil"
)

var prResting = "Вы присели и дали отдохнуть своим уставшим косточкам."
var prRestingAlt = "Вы удобно расположились на стуле у камина и начали отдыхать."
var prRestingAlter = "Вы дали отдохнуть своим уставшим косточкам."
var prSat = "Вы присели."
var prSatAlt = "Вы прекратили отдыхать и присели."
var prStoppedResting = "Вы прекратили отдыхать и поднялись на ноги."
var prStoodUp = "Вы встали на ноги."
var prWentToSleep = "Вы легли и заснули."
var prWentToSleepAlt = "Вы легли и заснули."
var prWokeUp = "Вы проснулись."
var prWokeUpAlt = "Вы потянулись, разминая затекшие мускулы и стряхивая остатки сладкого сна."

func (b *Bot) ParsePosition(s string) Event {
	if botutil.HasAnyLinePrefix(s, []string{prResting, prRestingAlt, prRestingAlter}) {
		b.Char.Position = POSITION_RESTING
		return EVENT_RESTING
	} else if botutil.HasAnyLinePrefix(s, []string{prSat, prSatAlt}) {
		b.Char.Position = POSITION_SITTING
		return EVENT_SAT
	} else if botutil.HasAnyLinePrefix(s, []string{prStoppedResting, prStoodUp}) {
		b.Char.Position = POSITION_STANDING
		return EVENT_STOOD_UP
	} else if botutil.HasAnyLinePrefix(s, []string{prWentToSleep, prWentToSleepAlt}) {
		b.Char.Position = POSITION_SLEEPING
		return EVENT_WENT_TO_SLEEP
	} else if botutil.HasAnyLinePrefix(s, []string{prWokeUp, prWokeUpAlt}) {
		b.Char.Position = POSITION_SITTING
		return EVENT_WOKE_UP
	} else {
		return EVENT_NOP
	}
}
