package repositories

import (
	"time"
)

// URLRepository определяет интерфейс для работы с URL
type URLRepository interface {
	SaveURL(longLink, shortLink string, expirationTime time.Time) error
	GetLongURL(shortLink string) (string, error)
	DeleteExpiredLinks() error
}
