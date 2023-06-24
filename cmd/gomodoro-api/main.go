package main

import (
	db "github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	srv "github.com/a-dakani/gomodoro/cmd/gomodoro-api/server"
	ws "github.com/a-dakani/gomodoro/cmd/gomodoro-api/ws"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not set")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("HOST is not set")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("PORT is not set")
	}

	// Start the websocket loop
	ws.Start()

	// Connect to the DB
	if err := db.ConnectDB(dsn); err != nil {
		log.Fatal(err)
	}

	// Start serving the API
	if err := srv.SetupAndListen(host, port); err != nil {
		log.Fatal(err)
	}
}
