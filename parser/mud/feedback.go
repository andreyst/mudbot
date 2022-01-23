package mud

import "mudbot/botutil"

var prCannotMoveInThisDirection = "Увы, Вы не можете идти в этом направлении."
var prCannotBecauseResting = "Вы слишком расслаблены, и Вам сейчас не до этого."

func (p *Parser) ParseFeedback(s string) (res []Event) {
	if botutil.HasLinePrefix(s, prCannotMoveInThisDirection) {
		res = append(res, EVENT_CANNOT_MOVE_IN_THIS_DIRECTION)
	}
	if botutil.HasLinePrefix(s, prCannotBecauseResting) {
		res = append(res, EVENT_CANNOT_BECAUSE_RESTING)
	}

	return
}
