package repositories

import (
	"time"
)

type URLRepository interface {
	SaveURL(longLink, shortLink string, expirationTime time.Time) error
	GetLongURL(shortLink string) (string, error)
	DeleteExpiredLinks() error
}
