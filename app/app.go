package app

import (
	"fmt"
	"mudbot/bot"
	"mudbot/proxy"
	"time"
)

type App struct {
	bot    *bot.Bot
	server *proxy.Server
}

func NewApp(localAddr string, remoteAddr string) *App {
	bot := bot.NewBot()
	server := proxy.NewServer(localAddr, remoteAddr, bot.Parse)

	app := App{
		bot:    bot,
		server: server,
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
