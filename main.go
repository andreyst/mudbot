package main

import (
	"fmt"
	"log"
	"mudbot/proxy"
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

	p := proxy.NewServer()
	p.Start(os.Args[1], os.Args[2])

	os.Exit(0)
}
