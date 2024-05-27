package postgres

import (
	"database/sql"
	"time"
)

type URLRepository struct {
	DB *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (repo *URLRepository) SaveURL(longLink, shortLink string, expirationTime time.Time) error {
	query := "INSERT INTO urls (long_link, short_link, expiration_time) VALUES ($1, $2, $3)"
	_, err := repo.DB.Exec(query, longLink, shortLink, expirationTime)
	return err
}

func (repo *URLRepository) GetLongURL(shortLink string) (string, error) {
	var longLink string
	err := repo.DB.QueryRow("SELECT long_link FROM urls WHERE short_link = $1", shortLink).Scan(&longLink)
	return longLink, err
}

func (repo *URLRepository) DeleteExpiredLinks() error {
	_, err := repo.DB.Exec("DELETE FROM urls WHERE expiration_time < NOW()")
	return err
}
