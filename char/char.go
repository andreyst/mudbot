package char

import (
	"mudbot/botutil"
	"strconv"
	"strings"

	"github.com/oriser/regroup"
	"go.uber.org/zap"
)

type Gender int

const (
	GENDER_MALE Gender = iota + 1
	GENDER_FEMALE
)

type Class int

const (
	CLASS_MERCENARY Class = iota + 1
	CLASS_RANGER
	CLASS_KNIGHT_OF_SOLAMNIA
	CLASS_THIEF
	CLASS_ASSASSIN
	CLASS_PRIEST_OF_MISHAKAL
	CLASS_PRIEST_OF_MORGION
	CLASS_MAGE
)

type MoralAlignment int

const (
	MORAL_ALIGNMENT_EVIL MoralAlignment = iota + 1
	MORAL_ALIGNMENT_NEUTRAL
	MORAL_ALIGNMENT_GOOD
)

type EthicalAlignment int

const (
	ETHICAL_ALIGNMENT_CHAOTIC EthicalAlignment = iota + 1
	ETHICAL_ALIGNMENT_NEUTRAL
	ETHICAL_ALIGNMENT_GOOD
)

type Position int

const (
	POSITION_STANDING Position = iota + 1
	POSITION_SITTING
	POSITION_RESTING
	POSITION_SLEEPING
)

type Char struct {
	logger *zap.SugaredLogger

	Stats Stats
}

type Stats struct {
	Name   string
	Age    int
	Gender Gender

	Class Class

	Health     int
	MaxHealth  int
	Stamina    int
	MaxStamina int

	IsHungry  bool
	IsThirsty bool

	MoralAlignment   MoralAlignment
	EthicalAlignment EthicalAlignment

	Level         int
	Experience    int
	ExperienceTNL int // Experience to next level

	Money int

	Position Position
}

var scoreMatcher = regroup.MustCompile(strings.Trim(`
(?ms)^Вы (?P<Gender>\pL+) (?P<Name>[^,]+), (?P<Class>\pL+) (?P<Level>\d+) уров\pL*\. Вам (?P<Age>\d+) \pL+\.$
У Вас сейчас (?P<Health>\d+) из (?P<MaxHealth>\d+) очк\pL+ жизни и (?P<Stamina>\d+) из (?P<MaxStamina>\d+) очк\pL+ бодрости.
.*
^Вы набрали (?P<Experience>\d+) очк\pL* опыта. До следующего уровня осталось (?P<ExperienceTNL>\d+)\.$
^У Вас есть (?P<Money>\d+) сталь\pL+ мон\pL+.$
.*
^Вы (?P<Position>(стоите|сидите|отдыхаете|спите))\.$`, "\n"))

func NewChar() *Char {
	c := Char{}
	c.logger = botutil.NewLogger("char")

	return &c
}

func (c *Char) ParseScore(s string) bool {
	// TODO: Add proper error handling
	match, _ := scoreMatcher.Groups(s)
	if match == nil {
		return false
	}

	c.Stats.Name = match["Name"]

	// TODO: handle conv errors
	c.Stats.Age, _ = strconv.Atoi(match["Age"])
	c.Stats.Level, _ = strconv.Atoi(match["Level"])
	c.Stats.Experience, _ = strconv.Atoi(match["Experience"])
	c.Stats.ExperienceTNL, _ = strconv.Atoi(match["ExperienceTNL"])
	c.Stats.Health, _ = strconv.Atoi(match["Health"])
	c.Stats.MaxHealth, _ = strconv.Atoi(match["MaxHealth"])
	c.Stats.Stamina, _ = strconv.Atoi(match["Stamina"])
	c.Stats.MaxStamina, _ = strconv.Atoi(match["MaxStamina"])
	c.Stats.Money, _ = strconv.Atoi(match["Money"])

	switch match["Class"] {
	case "наемник":
		c.Stats.Class = CLASS_MERCENARY
	case "следопыт":
		c.Stats.Class = CLASS_RANGER
	case "соламнийский рыцарь":
		c.Stats.Class = CLASS_KNIGHT_OF_SOLAMNIA
	case "вор":
		c.Stats.Class = CLASS_THIEF
	case "ассасин":
		c.Stats.Class = CLASS_ASSASSIN
	case "жрец Мишакаль":
		c.Stats.Class = CLASS_PRIEST_OF_MISHAKAL
	case "жрец Моргиона":
		c.Stats.Class = CLASS_PRIEST_OF_MORGION
	case "маг":
		c.Stats.Class = CLASS_MAGE
	default:
		c.logger.Warnf("Unknown class %+v", match["Class"])
		// TODO: Show to client warning about unknown class & log it
	}

	switch match["Gender"] {
	case "мужчина":
		c.Stats.Gender = GENDER_MALE
	case "женщина":
		c.Stats.Gender = GENDER_FEMALE
	default:
		c.logger.Warnf("Unknown gender %+v", match["Gender"])
		// TODO: Show to client warning about unknown gender & log it
	}

	switch match["Position"] {
	case "стоите":
		c.Stats.Position = POSITION_STANDING
	case "отдыхаете":
		c.Stats.Position = POSITION_RESTING
	case "сидите":
		c.Stats.Position = POSITION_SITTING
	case "спите":
		c.Stats.Position = POSITION_SLEEPING
	default:
		c.logger.Warnf("Unknown position %+v", match["Position"])
		// TODO: Show to client warning about unknown position & log it
	}

	if strings.Contains(s, "Ваша душа чиста, и Ваша вера в идеалы Добра непоколебима.") {
		c.Stats.EthicalAlignment = ETHICAL_ALIGNMENT_GOOD
	}
	if strings.Contains(s, "Вы привержены идеалам Добра.") {
		c.Stats.MoralAlignment = MORAL_ALIGNMENT_GOOD
	}

	c.Stats.IsHungry = strings.Contains(s, "Вы хотите есть.")
	c.Stats.IsThirsty = strings.Contains(s, "Вы хотите пить.")

	c.logger.Debugf("Char after score parse:\n%+v", c)

	return true
}

//go:generate stringer -type=Gender,Class,MoralAlignment,EthicalAlignment,Position -output char_string.go
