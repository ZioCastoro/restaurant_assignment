package services

import "github.com/ZioCastoro/restaurant_assignment/app/entities"

type AuthServiceInterface interface {
	GetAuthorizedUser(token entities.Token) (entities.User, error)
}
