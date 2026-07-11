package deps

import (
	"context"

	"github.com/ZioCastoro/restaurant_assignment/app"
	"github.com/ZioCastoro/restaurant_assignment/app/handlers"
	"github.com/ZioCastoro/restaurant_assignment/app/repositories"
	"github.com/ZioCastoro/restaurant_assignment/app/services"
	"github.com/ZioCastoro/restaurant_assignment/postgre"
)

func InjectApplication() (*app.Application, error) {
	dsn := "postgres://samu:samu@localhost:5432/samu?sslmode=disable"

	db, err := postgre.Connect(context.Background(), "postgres", dsn)
	if err != nil {
		return nil, err
	}

	expensesRepository := repositories.NewExpensesRepositoryPostgres(db)
	expensesService := services.NewExpensesService(expensesRepository)
	expenseHandler := handlers.NewExpensesHandler(expensesService)

	authService := services.NewAuthService()
	middleware := handlers.NewMiddleware(authService)

	return app.NewApplication(middleware, expenseHandler), nil
}
