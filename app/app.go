package app

import (
	"fmt"
	"mudbot/bot"
	"mudbot/botutil"
	"mudbot/proxy"
	"os"
	"time"

	"go.uber.org/zap"
)

type App struct {
	bot    *bot.Bot
	server *proxy.Server

	logger *zap.SugaredLogger
}

func NewApp(localAddr string, remoteAddr string) *App {
	logger := botutil.NewLogger("app")

	login, hasLogin := os.LookupEnv("LOGIN")
	if !hasLogin {
		logger.Fatalf("LOGIN env var missing")
	}
	password, hasPassword := os.LookupEnv("PASSWORD")
	if !hasPassword {
		logger.Fatalf("PASSWORD env var missing")
	}

	bot := bot.NewBot(bot.Credentials{
		Login:    login,
		Password: password,
	})

	server := proxy.NewServer(localAddr, remoteAddr, bot.Parse)
	bot.SetToMudSender(server.SendToMud)

	app := App{
		bot:    bot,
		server: server,
		logger: logger,
	}

	return &app
}

func (app *App) Start() {
	go func() {
		for {
			fmt.Printf("Bot: %+v\n", app.bot)
			time.Sleep(time.Duration(2) * time.Second)
		}
	}()

	app.server.Start()
}
