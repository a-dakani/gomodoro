package server

import (
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

// TODO a more specific Error Handling without sending the whole error back to the client

func getTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("gomodoro_id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing id. Must be an integer " + err.Error(),
		})
	}

	gomodoroID32 := uint(gomodoroID64)

	timer, err := model.GetTimer(gomodoroID32)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting latest timer" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}

func startTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("gomodoro_id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing id. Must be an integer " + err.Error(),
		})
	}

	gomodoroID32 := uint(gomodoroID64)

	timer, err := model.GetTimer(gomodoroID32)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting timer " + err.Error(),
		})
	}

	timer.Status = model.Running
	timer.StartedAt = time.Now()

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)

}

func pauseTimer(ctx *fiber.Ctx) error {
	return nil
}

func resumeTimer(ctx *fiber.Ctx) error {
	return nil
}

func resetTimer(ctx *fiber.Ctx) error {
	return nil
}
func stopTimer(ctx *fiber.Ctx) error {
	return nil
}

func nextTimer(ctx *fiber.Ctx) error {
	return nil
}
