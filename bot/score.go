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

func (b *Bot) ParseScore(s string) Event {
	// TODO: Add proper error handling
	match, _ := scoreMatcher.Groups(s)
	if match == nil {
		return EVENT_NOP
	}

	b.Char.Initialized = true

	b.Char.Name = match["Name"]

	// TODO: handle conv errors
	b.Char.Age, _ = strconv.Atoi(match["Age"])
	b.Char.Level, _ = strconv.Atoi(match["Level"])
	b.Char.Experience, _ = strconv.Atoi(match["Experience"])
	b.Char.ExperienceTNL, _ = strconv.Atoi(match["ExperienceTNL"])
	b.Char.Health, _ = strconv.Atoi(match["Health"])
	b.Char.MaxHealth, _ = strconv.Atoi(match["MaxHealth"])
	b.Char.Stamina, _ = strconv.Atoi(match["Stamina"])
	b.Char.MaxStamina, _ = strconv.Atoi(match["MaxStamina"])
	b.Char.Money, _ = strconv.Atoi(match["Money"])

	switch match["Class"] {
	case "наемник":
		b.Char.Class = CLASS_MERCENARY
	case "следопыт":
		b.Char.Class = CLASS_RANGER
	case "соламнийский рыцарь":
		b.Char.Class = CLASS_KNIGHT_OF_SOLAMNIA
	case "вор":
		b.Char.Class = CLASS_THIEF
	case "ассасин":
		b.Char.Class = CLASS_ASSASSIN
	case "жрец Мишакаль":
		b.Char.Class = CLASS_PRIEST_OF_MISHAKAL
	case "жрец Моргиона":
		b.Char.Class = CLASS_PRIEST_OF_MORGION
	case "маг":
		b.Char.Class = CLASS_MAGE
	default:
		b.WarnClientf("Unknown class %v", match["Class"])
		// TODO: Log about unknown class
	}

	switch match["Gender"] {
	case "мужчина":
		b.Char.Gender = GENDER_MALE
	case "женщина":
		b.Char.Gender = GENDER_FEMALE
	default:
		b.WarnClientf("Unknown gender %v", match["Gender"])
		// TODO: Log it
	}

	switch match["Position"] {
	case "стоите":
		b.Char.Position = POSITION_STANDING
	case "отдыхаете":
		b.Char.Position = POSITION_RESTING
	case "сидите":
		b.Char.Position = POSITION_SITTING
	case "спите":
		b.Char.Position = POSITION_SLEEPING
	default:
		b.WarnClientf("Unknown position %+v", match["Position"])
		// TODO: Log it
	}

	if strings.Contains(s, "Ваша душа чиста, и Ваша вера в идеалы Добра непоколебима.") {
		b.Char.EthicalAlignment = ETHICAL_ALIGNMENT_GOOD
	}
	if strings.Contains(s, "Вы привержены идеалам Добра.") {
		b.Char.MoralAlignment = MORAL_ALIGNMENT_GOOD
	}

	b.logger.Debugf("Char after score parse:\n%+v", b.Char)

	return EVENT_SCORE
}
