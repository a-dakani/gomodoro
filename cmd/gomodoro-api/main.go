package main

import (
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/server"
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

	if err := model.ConnectDB(dsn); err != nil {
		log.Fatal(err)
	}

	if err := server.SetupAndListen(host, port); err != nil {
		log.Fatal(err)
	}
}
