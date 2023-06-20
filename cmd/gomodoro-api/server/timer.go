package server

import (
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

// TODO a more specific Error Handling without sending the whole error back to the client
// TODO don't allow to start a timer if there is already one running

func getTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
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
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
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
	timer.Remaining = timer.Duration

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}

func pauseTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
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

	timer.Status = model.Paused
	timer.Remaining = timer.Duration - time.Since(timer.StartedAt)

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}

func resumeTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
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
	timer.StartedAt = time.Now().Add(-timer.Remaining)

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}

func resetTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
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

	timer.Status = model.Idle
	timer.StartedAt = time.Time{}
	timer.Remaining = timer.Duration

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}

func nextTimer(ctx *fiber.Ctx) error {
	gomodoroID64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing id. Must be an integer " + err.Error(),
		})
	}

	gomodoroID32 := uint(gomodoroID64)

	gomodoro, err := model.GetGomodoroByID(gomodoroID32)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting gomodoro " + err.Error(),
		})
	}

	timer, err := model.GetTimer(gomodoroID32)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting timer " + err.Error(),
		})
	}

	if timer.Type == model.WorkTimer {
		if timer.Repetition == gomodoro.Repetitions {
			timer.Type = model.LongBreakTimer
			timer.Duration = gomodoro.LongBreak
			timer.Remaining = gomodoro.LongBreak
		} else {
			timer.Type = model.ShortBreakTimer
			timer.Duration = gomodoro.ShortBreak
			timer.Remaining = gomodoro.ShortBreak
		}

	} else if timer.Type == model.ShortBreakTimer {
		timer.Type = model.WorkTimer
		timer.Duration = gomodoro.Work
		timer.Remaining = gomodoro.Work
		timer.Repetition++

	} else if timer.Type == model.LongBreakTimer {
		timer.Type = model.WorkTimer
		timer.Duration = gomodoro.Work
		timer.Remaining = gomodoro.Work
		timer.Repetition = 1

	}

	err = model.UpdateTimer(timer.ID, &timer)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(timer)
}
