package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func SetupAndListen() error {
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
	apiV1.Post("/g/c/:name", createGomodoro)
	apiV1.Delete("/g/:id", deleteGomodoroByID)
	apiV1.Delete("/g/n/:name", deleteGomodoroByName)
	apiV1.Put("/g/:id", updateGomodoro)

	apiV1.Get("/g/:id/t/", getTimer)
	apiV1.Put("/g/:id/t/start", startTimer)
	apiV1.Put("/g/:id/t/stop", stopTimer)
	apiV1.Put("/g/:id/t/pause", pauseTimer)
	apiV1.Put("/g/:id/t/resume", resumeTimer)
	apiV1.Put("/g/:id/t/reset", resetTimer)
	apiV1.Put("/g/:id/t/next", nextTimer)

	router.Use("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Method Not Found",
		})
	})

	if err := router.Listen("localhost:3000"); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
