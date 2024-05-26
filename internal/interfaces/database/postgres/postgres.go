package postgres

import (
	"URL_SHORT/pkg/config"
	"database/sql"
	"fmt"
)

func ConnectToDB(cfg config.DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Successfully connected to PostgreSQL")
	return db, nil
}
