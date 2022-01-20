package bot

import (
	"mudbot/botutil"
	"strings"

	"go.uber.org/zap"
)

type Bot struct {
	char  Char
	fight Fight

	logger *zap.SugaredLogger
}

func NewBot() *Bot {
	bot := Bot{
		logger: botutil.NewLogger("bot"),
	}

	return &bot
}

func (b *Bot) Parse(chunk []byte) {
	s := strings.ReplaceAll(string(chunk), "\r\n", "\n")

	b.ParseScore(s)
	b.ParsePrompt(s)
}
