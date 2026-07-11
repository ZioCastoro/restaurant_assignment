package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDKey = "user_id"

var ErrNotAuthorized = errors.New("not authorized")

type Middleware struct {
	authService services.AuthServiceInterface
}

func (m *Middleware) AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := uuid.Parse(bearerToken)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	user, err := m.authService.GetAuthorizedUser(entities.Token(token))
	if err != nil {
		return fmt.Errorf("not authorized: %w", err)
	}

	ctx.Locals(userIDKey, user.ID.String())

	return ctx.Next()
}

func NewMiddleware(
	authService services.AuthServiceInterface,
) *Middleware {
	return &Middleware{
		authService: authService,
	}
}
