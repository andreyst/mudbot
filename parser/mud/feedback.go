package mud

import "mudbot/botutil"

var prCannotMoveInThisDirection = "Увы, Вы не можете идти в этом направлении."
var sufCannotBecauseClosedMasculine = "закрыт."
var sufCannotBecauseClosedFeminine = "закрыта."
var sufCannotBecauseClosedNeutral = "закрыто."
var sufCannotBecauseClosedPlural = "закрыты."
var prCannotBecauseResting = "Вы слишком расслаблены, и Вам сейчас не до этого."
var prCannotBecauseSitting = "Для этого необходимо встать на ноги."
var prCannotBecauseSleeping = "Сделать это во сне, что ли?"

func (p *Parser) ParseFeedback(s string) (res []Event) {
	if botutil.HasLinePrefix(s, prCannotMoveInThisDirection) {
		res = append(res, EVENT_CANNOT_MOVE_IN_THIS_DIRECTION)
	}
	if botutil.HasLinePrefix(s, prCannotBecauseResting) {
		res = append(res, EVENT_CANNOT_BECAUSE_RESTING)
	}
	if botutil.HasLinePrefix(s, prCannotBecauseSitting) {
		res = append(res, EVENT_CANNOT_BECAUSE_SITTING)
	}
	if botutil.HasLinePrefix(s, prCannotBecauseSleeping) {
		res = append(res, EVENT_CANNOT_BECAUSE_SLEEPING)
	}
	if botutil.HasAnyLineSuffix(s, []string{
		sufCannotBecauseClosedMasculine,
		sufCannotBecauseClosedFeminine,
		sufCannotBecauseClosedNeutral,
		sufCannotBecauseClosedPlural,
	}) {
		res = append(res, EVENT_CANNOT_BECAUSE_CLOSED)
	}

	return
}
