package handlers

import (
	"github.com/ZioCastoro/restaurant_assignment/app/handlers/payloads"
	"github.com/ZioCastoro/restaurant_assignment/app/services"
	"github.com/gofiber/fiber/v2"
)

const paramIDKey = "id"

type ExpensesHandler struct {
	service services.ExpenseServiceInterface
}

func (e *ExpensesHandler) List(ctx *fiber.Ctx) error {
	var req payloads.ExpenseFilterRequest
	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	filters := req.MapToEntity()

	expenses, err := e.service.List(ctx.Context(), filters)
	if err != nil {
		return err
	}

	res := payloads.NewPaginatedExpensesResponse(expenses, filters.Pagination)

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (e *ExpensesHandler) Find(ctx *fiber.Ctx) error {
	id := ctx.Params(paramIDKey)

	expense, err := e.service.Find(ctx.Context(), id)
	if err != nil {
		return err
	}

	res := payloads.NewExpenseResponse(expense)

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (e *ExpensesHandler) Create(ctx *fiber.Ctx) error {
	var req payloads.CreateExpenseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	expense := req.MapToEntity()

	expense, err := e.service.Create(ctx.Context(), expense)
	if err != nil {
		return err
	}

	res := payloads.NewExpenseResponse(expense)

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (e *ExpensesHandler) Update(ctx *fiber.Ctx) error {
	var req payloads.UpdateExpenseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	expense := req.MapToEntity(ctx.Params(paramIDKey))

	expense, err := e.service.Update(ctx.Context(), expense)
	if err != nil {
		return err
	}

	res := payloads.NewExpenseResponse(expense)

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (e *ExpensesHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params(paramIDKey)

	if err := e.service.Delete(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (e *ExpensesHandler) Import(ctx *fiber.Ctx) error {
	var req payloads.ImportExpenseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	expenses := req.MapToEntity()

	go e.service.Import(ctx.Context(), expenses)

	return ctx.SendStatus(fiber.StatusAccepted)
}

func NewExpensesHandler(
	service services.ExpenseServiceInterface,
) *ExpensesHandler {
	return &ExpensesHandler{
		service: service,
	}
}
