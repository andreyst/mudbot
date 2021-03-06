package app

import (
	"encoding/json"
	"fmt"
	"mudbot/atlas"
	"mudbot/bot"
	"mudbot/botutil"
	"mudbot/parser/client"
	"mudbot/parser/mud"
	"mudbot/proxy"
	"os"
	"strconv"
	"strings"
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
	a.StartServer()

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
		atlas:  a,
		bot:    b,
		server: server,
		logger: logger,
	}

	return &app
}

func (app *App) Start() {
	go func() {
		botFo, err := os.Create("/tmp/mudbot")
		if err != nil {
			panic(err)
		}

		atlasFo, err := os.Create("/tmp/atlas")
		if err != nil {
			panic(err)
		}

		atlasDotFo, err := os.Create("/tmp/atlas.dot")
		if err != nil {
			panic(err)
		}

		for {
			_, botSeekErr := botFo.Seek(0, 0)
			if botSeekErr != nil {
				panic(err)
			}
			_, botWriteErr := botFo.Write([]byte(fmt.Sprintf("%+v", app.bot)))
			if botWriteErr != nil {
				panic(err)
			}

			_, atlasSeekErr := atlasFo.Seek(0, 0)
			if atlasSeekErr != nil {
				panic(err)
			}
			atlasJson, atlasMarshalErr := json.Marshal(app.atlas)
			if atlasMarshalErr != nil {
				panic(err)
			}
			_, atlasWriteErr := atlasFo.Write(atlasJson)
			if atlasWriteErr != nil {
				panic(err)
			}

			_, atlasDotSeekErr := atlasDotFo.Seek(0, 0)
			if atlasDotSeekErr != nil {
				panic(err)
			}
			var dot strings.Builder
			dot.WriteString("digraph G {\n")
			for roomId, room := range app.atlas.Rooms {
				roomIdStr := strconv.FormatInt(roomId, 10)
				dot.WriteString("  r" + roomIdStr + " [label=\"" + strings.ReplaceAll(room.Name, "\"", "\\\"") + "\"]\n")
				for exitDir, exitRoomId := range room.Exits {
					exitRoomIdStr := strconv.FormatInt(exitRoomId, 10)
					dot.WriteString("  r" + roomIdStr + " -> r" + exitRoomIdStr + " [ label=\"" + exitDir.String() + "\" ];\n")
				}
				dot.WriteString("\n")
			}
			dot.WriteString("}\n")
			_, atlasDotWriteErr := atlasDotFo.Write([]byte(dot.String()))
			if atlasDotWriteErr != nil {
				panic(err)
			}

			time.Sleep(time.Duration(2) * time.Second)
		}
	}()

	app.server.Start()
}
