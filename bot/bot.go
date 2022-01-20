package bot

import (
	"mudbot/botutil"
	"strings"

	"go.uber.org/zap"
)

type Sender func(string, bool)

type Bot struct {
	Credentials Credentials
	Char        Char
	Fight       Fight

	InGame bool

	toMudSender Sender
	logger      *zap.SugaredLogger
}

func NewBot(credentials Credentials) *Bot {
	bot := Bot{
		Credentials: credentials,
		logger:      botutil.NewLogger("bot"),
	}

	return &bot
}

func (b *Bot) SetToMudSender(f Sender) {
	b.toMudSender = func(s string, echo bool) { f(s+"\r\n", echo) }
}

func (b *Bot) SendToMud(s string) {
	b.toMudSender(s, true)
}

func (b *Bot) SendToMudWithoutEcho(s string) {
	b.toMudSender(s, false)
}

func (b *Bot) Parse(chunk []byte) {
	s := strings.ReplaceAll(string(chunk), "\r\n", "\n")

	if !b.InGame {
		b.ParseLogin(s)
	} else {
		b.ParseScore(s)
		b.ParsePrompt(s)
	}
}

func (b *Bot) Step() {
	if !b.InGame {
		return
	}

	if !b.Fight.IsActive {
		if !b.Char.Initialized {
			b.SendToMud("score")
		}
	}
}
