package repositories

import (
	"URL_SHORT/internal/domain/models"
	_ "URL_SHORT/internal/domain/models"
)

type UserRepository interface {
	SaveUser(u models.User) error
	GetUser(u models.User) error
	DeleteUser(u models.User) error
	UpdateUser(u models.User) error
}
