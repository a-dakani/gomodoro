package server

import (
	"github.com/a-dakani/gomodoro/cmd/gomodoro-api/model"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// TODO a more specific Error Handling without sending the whole error back to the client

func getAllGomodoros(ctx *fiber.Ctx) error {
	gomodoros, err := model.GetAllGomodoros()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting all gomodoros " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(gomodoros)
}

func getGomodoroByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	gomodoro, err := model.GetGomodoroByName(name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting gomodoro " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(gomodoro)
}

func getGomodoroByID(ctx *fiber.Ctx) error {
	id64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing id. Must be an integer " + err.Error(),
		})
	}

	id32 := uint(id64)

	gomodoro, err := model.GetGomodoroByID(id32)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Error getting gomodoro " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(gomodoro)
}

func createGomodoro(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	gomodoro, err := model.CreateGomodoro(name)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error creating gomodoro " + err.Error(),
		})
	}

	err = model.CreateDefaultTimer(gomodoro.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating gomodoro timer " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(gomodoro)
}

func deleteGomodoroByID(ctx *fiber.Ctx) error {
	id64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing id. Must be an integer " + err.Error(),
		})
	}

	id32 := uint(id64)

	if err := model.DeleteGomodoroByID(id32); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting gomodoro" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func deleteGomodoroByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	if err := model.DeleteGomodoroByName(name); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting gomodoro" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func updateGomodoro(ctx *fiber.Ctx) error {
	id64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error parsing id" + err.Error(),
		})
	}

	id32 := uint(id64)

	gomodoro := new(model.Gomodoro)

	if err := ctx.BodyParser(gomodoro); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing gomodoro" + err.Error(),
		})
	}

	if err := model.UpdateGomodoro(id32, gomodoro); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating gomodoro" + err.Error(),
		})
	}

	newGomodoro, err := model.GetGomodoroByID(id32)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting gomodoro" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(newGomodoro)
}
