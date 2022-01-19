//go:generate stringer -type=Gender,Class,MoralAlignment,EthicalAlignment,Position,Condition -output types_string.go

package bot

type Gender int

const (
	GENDER_UNKNOWN Gender = iota
	GENDER_MALE
	GENDER_FEMALE
)

type Class int

const (
	CLASS_UNKNOWN Class = iota
	CLASS_MERCENARY
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
	MORAL_ALIGNMENT_UNKNOWN MoralAlignment = iota
	MORAL_ALIGNMENT_EVIL
	MORAL_ALIGNMENT_NEUTRAL
	MORAL_ALIGNMENT_GOOD
)

type EthicalAlignment int

const (
	ETHICAL_ALIGNMENT_UNKNOWN EthicalAlignment = iota
	ETHICAL_ALIGNMENT_CHAOTIC
	ETHICAL_ALIGNMENT_NEUTRAL
	ETHICAL_ALIGNMENT_GOOD
)

type Position int

const (
	POSITION_UNKNOWN Position = iota
	POSITION_STANDING
	POSITION_SITTING
	POSITION_LYING
	POSITION_RESTING
	POSITION_SLEEPING
)

type Char struct {
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
	ExperienceTNL int // TNL = to next level

	Money int

	Position Position
}

type Fight struct {
	IsActive bool
	HaveTank bool
	Tank     Combatant
	Enemy    Combatant
}

type Combatant struct {
	Name      string
	Condition Condition
	Position  Position
}

type Condition int

const (
	CONDITION_UNKNOWN Condition = iota
	CONDITION_PERFECT
	CONDITION_VERY_GOOD
	CONDITION_GOOD
	CONDITION_AVERAGE
	CONDITION_BAD
	CONDITION_VERY_BAD
	CONDITION_AWFUL
	CONDITION_DYING
)
