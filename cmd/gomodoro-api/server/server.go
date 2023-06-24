package server

import (
	"fmt"
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func SetupAndListen(host string, port int) error {
	router := fiber.New()
	router.Use(cors.New(cors.Config{
		// TODO: Change this to the frontend URL
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Use(logger.New(logger.Config{}))

	apiV1 := router.Group("/api/v1")

	apiV1.Get("/g/", getAllGomodoros)
	apiV1.Get("/g/:id", getGomodoroByID)
	apiV1.Get("/g/n/:name", getGomodoroByName)
	apiV1.Post("/g/n/:name", createGomodoro)
	apiV1.Delete("/g/:id", deleteGomodoroByID)
	apiV1.Delete("/g/n/:name", deleteGomodoroByName)
	apiV1.Put("/g/:id", updateGomodoro)

	apiV1.Get("/g/:id/t/", getTimer)
	apiV1.Put("/g/:id/t/start", startTimer)
	apiV1.Put("/g/:id/t/reset", resetTimer)
	apiV1.Put("/g/:id/t/pause", pauseTimer)
	apiV1.Put("/g/:id/t/resume", resumeTimer)
	apiV1.Put("/g/:id/t/next", nextTimer)

	apiV1.Get("/g/:id/ws", websocket.New(ws.Serve))

	router.Use("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Method Not Found",
		})
	})

	log.Printf("Listening on %s:%d", host, port)

	if err := router.Listen(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
