package app

import (
	"github.com/ZioCastoro/restaurant_assignment/app/handlers"
	"github.com/gofiber/fiber/v2"
)

type Application struct {
	middleware      handlers.Middleware
	expensesHandler handlers.ExpensesHandler
}

func (a *Application) Routes() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: errorHandler,
			UnescapePath: true,
			BodyLimit:    64 * 1024 * 1024,
		},
	)

	expenses := app.Group("/expenses")

	expenses.Use(a.middleware.AuthMiddleware)

	expenses.Get("/", a.expensesHandler.List)
	expenses.Get("/:id/", a.expensesHandler.Find)
	expenses.Post("/", a.expensesHandler.Create)
	expenses.Put("/:id/", a.expensesHandler.Update)
	expenses.Delete("/:id/", a.expensesHandler.Delete)
	expenses.Post("/import/", a.expensesHandler.Import)

	// fallback to 404.
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})

	return app
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	var res struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}

	res.Status = fiber.StatusBadRequest
	res.Error = err.Error()

	// TODO: Handle different error codes.
	return ctx.Status(fiber.StatusBadRequest).JSON(res)
}

func NewApplication(
	middleware *handlers.Middleware,
	expensesHandler *handlers.ExpensesHandler,
) *Application {
	return &Application{
		middleware:      *middleware,
		expensesHandler: *expensesHandler,
	}
}
