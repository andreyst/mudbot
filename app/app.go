package app

import (
	"fmt"
	"mudbot/atlas"
	"mudbot/bot"
	"mudbot/botutil"
	"mudbot/parser/client"
	"mudbot/parser/mud"
	"mudbot/proxy"
	"os"
	"time"

	"go.uber.org/zap"
)

type App struct {
	bot    *bot.Bot
	server *proxy.Server
	atlas  *atlas.Atlas

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

	a := atlas.NewAtlas()

	b := bot.NewBot(bot.Credentials{
		Login:    login,
		Password: password,
	})

	clientParser := client.NewParser(a)
	mudParser := mud.NewParser(a)

	server := proxy.NewServer(localAddr, remoteAddr, clientParser.Parse, func(bytes []byte) {
		mudParser.Parse(bytes)
		b.Parse(bytes)
	})
	b.SetToMudSender(server.SendToMud)
	b.SetToClientSender(server.SendToClient)

	app := App{
		bot:    b,
		server: server,
		logger: logger,
	}

	return &app
}

func (app *App) Start() {
	go func() {
		fo, err := os.Create("/tmp/mudbot")
		if err != nil {
			panic(err)
		}
		for {
			_, err := fo.Seek(0, 0)
			_, err = fo.Write([]byte(fmt.Sprintf("%+v", app.bot)))
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(2) * time.Second)
		}
	}()

	app.server.Start()
}
