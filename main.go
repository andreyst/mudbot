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

	s := proxy.NewServer(os.Args[1], os.Args[2])
	s.Start()

	os.Exit(0)
}
