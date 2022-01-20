package main

import (
	"fmt"
	"log"
	"mudbot/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v <local host:port> <mud host:port>", os.Args[0])
		os.Exit(1)
	}

	app := app.NewApp(os.Args[1], os.Args[2])
	app.Start()

	os.Exit(0)
}
