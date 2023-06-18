package main

import (
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DSN")

	if err := model.ConnectDB(dsn); err != nil {
		log.Fatal(err)
	}

	if err := server.SetupAndListen(); err != nil {
		log.Fatal(err)
	}
}
