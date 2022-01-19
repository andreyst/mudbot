package bot

import (
	"strconv"
	"strings"

	"github.com/oriser/regroup"
)

var scoreMatcher = regroup.MustCompile(strings.Trim(`
(?ms)^Вы (?P<Gender>\pL+) (?P<Name>[^,]+), (?P<Class>\pL+) (?P<Level>\d+) уров\pL*\. Вам (?P<Age>\d+) \pL+\.$
У Вас сейчас (?P<Health>\d+) из (?P<MaxHealth>\d+) очк\pL+ жизни и (?P<Stamina>\d+) из (?P<MaxStamina>\d+) очк\pL+ бодрости.
.*
^Вы набрали (?P<Experience>\d+) очк\pL* опыта. До следующего уровня осталось (?P<ExperienceTNL>\d+)\.$
^У Вас есть (?P<Money>\d+) сталь\pL+ мон\pL+.$
.*
^Вы (?P<Position>(стоите|сидите|отдыхаете|спите))\.$`, "\n"))

func (b *Bot) ParseScore(s string) bool {
	// TODO: Add proper error handling
	match, _ := scoreMatcher.Groups(s)
	if match == nil {
		return false
	}

	b.char.Name = match["Name"]

	// TODO: handle conv errors
	b.char.Age, _ = strconv.Atoi(match["Age"])
	b.char.Level, _ = strconv.Atoi(match["Level"])
	b.char.Experience, _ = strconv.Atoi(match["Experience"])
	b.char.ExperienceTNL, _ = strconv.Atoi(match["ExperienceTNL"])
	b.char.Health, _ = strconv.Atoi(match["Health"])
	b.char.MaxHealth, _ = strconv.Atoi(match["MaxHealth"])
	b.char.Stamina, _ = strconv.Atoi(match["Stamina"])
	b.char.MaxStamina, _ = strconv.Atoi(match["MaxStamina"])
	b.char.Money, _ = strconv.Atoi(match["Money"])

	switch match["Class"] {
	case "наемник":
		b.char.Class = CLASS_MERCENARY
	case "следопыт":
		b.char.Class = CLASS_RANGER
	case "соламнийский рыцарь":
		b.char.Class = CLASS_KNIGHT_OF_SOLAMNIA
	case "вор":
		b.char.Class = CLASS_THIEF
	case "ассасин":
		b.char.Class = CLASS_ASSASSIN
	case "жрец Мишакаль":
		b.char.Class = CLASS_PRIEST_OF_MISHAKAL
	case "жрец Моргиона":
		b.char.Class = CLASS_PRIEST_OF_MORGION
	case "маг":
		b.char.Class = CLASS_MAGE
	default:
		b.logger.Warnf("Unknown class %+v", match["Class"])
		// TODO: Show to client warning about unknown class & log it
	}

	switch match["Gender"] {
	case "мужчина":
		b.char.Gender = GENDER_MALE
	case "женщина":
		b.char.Gender = GENDER_FEMALE
	default:
		b.logger.Warnf("Unknown gender %+v", match["Gender"])
		// TODO: Show to client warning about unknown gender & log it
	}

	switch match["Position"] {
	case "стоите":
		b.char.Position = POSITION_STANDING
	case "отдыхаете":
		b.char.Position = POSITION_RESTING
	case "сидите":
		b.char.Position = POSITION_SITTING
	case "спите":
		b.char.Position = POSITION_SLEEPING
	default:
		b.logger.Warnf("Unknown position %+v", match["Position"])
		// TODO: Show to client warning about unknown position & log it
	}

	if strings.Contains(s, "Ваша душа чиста, и Ваша вера в идеалы Добра непоколебима.") {
		b.char.EthicalAlignment = ETHICAL_ALIGNMENT_GOOD
	}
	if strings.Contains(s, "Вы привержены идеалам Добра.") {
		b.char.MoralAlignment = MORAL_ALIGNMENT_GOOD
	}

	b.char.IsHungry = strings.Contains(s, "Вы хотите есть.")
	b.char.IsThirsty = strings.Contains(s, "Вы хотите пить.")

	b.logger.Debugf("Char after score parse:\n%+v", b.char)

	return true
}
