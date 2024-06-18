package postgres

import (
	"URL_SHORT/internal/domain/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) SaveUser(u *models.User) error {
	query := `INSERT INTO users (username, password_hash, email, is_activated) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, created_at, is_activated, version`
	err := repo.DB.QueryRow(query, u.Username, u.Password, u.Email, u.IsActivated).Scan(&u.ID, &u.CreatedAt, &u.IsActivated, &u.Version)
	return err
}

func (repo *UserRepository) GetUser(email string, u *models.User) (*models.User, error) {
	query := "SELECT id, username, password_hash, created_at, email, is_activated, version FROM users WHERE email = $1"
	err := repo.DB.QueryRow(query, email).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.Email, &u.IsActivated, &u.Version)
	return u, err
}

func (repo *UserRepository) DeleteUser(email string) error {
	query := "DELETE FROM users WHERE email = $1"
	_, err := repo.DB.Exec(query, email)
	return err
}

func (repo *UserRepository) UpdateUser(u *models.User) (*models.User, error) {
	query := "UPDATE users SET username = $1, password_hash = $2, email = $3, is_activated = $4, version = $5 WHERE id = $6"
	_, err := repo.DB.Exec(query, u.Username, u.Password, u.Email, u.IsActivated, u.Version, u.ID)
	return u, err
}
