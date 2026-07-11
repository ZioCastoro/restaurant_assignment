package services

import (
	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/google/uuid"
)

type AuthService struct{}

func (a *AuthService) GetAuthorizedUser(token entities.Token) (entities.User, error) {
	// TODO: Implement real authentication logic.
	var mockUser = entities.User{
		ID:       uuid.New(),
		Username: "zio",
	}

	return mockUser, nil
}

// Guard for interface implementation.
var _ AuthServiceInterface = new(AuthService)

func NewAuthService() *AuthService {
	return &AuthService{}
}
