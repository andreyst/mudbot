package bot

import (
	"fmt"
	"mudbot/botutil"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Sender func(string)
type EchoSender func(string, bool)

type Bot struct {
	Credentials Credentials
	Char        Char
	Fight       Fight

	InGame bool

	State State

	toMudSender    EchoSender
	toClientSender Sender

	logger *zap.SugaredLogger
}

func NewBot(credentials Credentials) *Bot {
	bot := Bot{
		Credentials: credentials,
		logger:      botutil.NewLogger("bot"),
	}

	return &bot
}

func (b *Bot) SetToMudSender(f EchoSender) {
	b.toMudSender = func(s string, echo bool) { f(s+"\r\n", echo) }
}

func (b *Bot) SetToClientSender(f Sender) {
	b.toClientSender = func(s string) { f(s) }
}

func (b *Bot) SendToMud(s string) {
	b.toMudSender(s, true)
}

func (b *Bot) SendToMudWithoutEcho(s string) {
	b.toMudSender(s, false)
}

func (b *Bot) SendToClient(s string) {
	b.toClientSender(s)
}

func (b *Bot) WarnClient(s string) {
	b.logger.Warnf(s)

	c := color.New(color.FgHiYellow)
	c.EnableColor()
	b.toClientSender(c.SprintFunc()("WARN: " + s))
}

func (b *Bot) WarnClientf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	b.WarnClient(s)
}

func (b *Bot) ErrorClient(s string) {
	b.logger.Errorf(s)

	c := color.New(color.FgHiRed)
	c.EnableColor()
	b.toClientSender(c.SprintFunc()("ERROR: " + s))
}

func (b *Bot) ErrorClientf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	b.ErrorClient(s)
}

func (b *Bot) Parse(bytes []byte) {
	// TODO: Remove this in favor of splitting multiline regexp
	// to decrease level of implicitness ("who changed linefeeds?")
	s := strings.ReplaceAll(string(bytes), "\r\n", "\n")

	b.ProcessEvent(b.ParseLogin(s))
	b.ProcessEvent(b.ParseScore(s))
	b.ProcessEvent(b.ParsePosition(s))
	b.ProcessEvents(b.ParseConsumption(s))
	b.ProcessEvents(b.ParseFeedback(s))
	b.ProcessEvent(b.ParsePrompt(s)) // Should always go last to trigger prompt event with all other data
}
