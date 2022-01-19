package bot

import (
	"regexp"
	"strings"
)

var combatantRe = regexp.MustCompile(`\[([^\]]+)\]`)

func (b *Bot) ParseFight(s string) {
	combatantMatch := combatantRe.FindAllStringSubmatch(s, -1)

	b.fight.IsActive = len(combatantMatch) > 0
	b.fight.HaveTank = len(combatantMatch) > 2

	for combatantNum, combatantStr := range combatantMatch {
		combatant := b.ParseCombatant(combatantStr[1])

		switch combatantNum {
		case 0:
			if len(combatantMatch) == 1 {
				b.fight.Enemy = combatant
			} else {
				b.char.Position = combatant.Position
			}
		case 1:
			if len(combatantMatch) == 2 {
				b.fight.Enemy = combatant
			} else {
				b.fight.Tank = combatant
			}
		case 2:
			b.fight.Enemy = combatant
		default:
			// TODO: Warn client
			b.logger.Warnf("Unknown combatantNum: %v", combatantNum)
		}
	}
}
func (b *Bot) ParseCombatant(s string) Combatant {
	parts := strings.Split(s, ":")
	subParts := strings.Split(parts[1], ",")

	var condition Condition
	switch subParts[0] {
	case "великолепное":
		condition = CONDITION_PERFECT
	case "среднее":
		condition = CONDITION_AVERAGE
	case "умирает":
		condition = CONDITION_DYING
	default:
		// TODO: Warn to client
		b.logger.Warnf("Unknown condition: %v", subParts[0])
	}

	position := POSITION_STANDING
	if len(subParts) > 1 {
		switch subParts[1] {
		case "леж":
			position = POSITION_LYING
		default:
			// TODO: Warn to client
			b.logger.Warnf("Unknown position: %v", subParts[1])
		}
	}

	return Combatant{
		Name:      parts[0],
		Condition: condition,
		Position:  position,
	}
}
