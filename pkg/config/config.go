package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

func getMigrationsPath() (string, error) {
	absPath, err := filepath.Abs("./migrations")
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func RunMigrations() error {

	dbConfig := LoadDBConfig()
	connString := fmt.Sprintf("'postgres://%s:%s@%s:%s/%s?sslmode=%s'", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode)

	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return err
	}

	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("Подключение невозможно: %v", err)
	}
	defer db.Close()

	driver, err := postgres.
		fmt.Println(dsn)
	cmd := exec.Command("migrate", "-path", fmt.Sprintf("file:///%s", migrationsPath), "-database", "pgx://"+dsn, "up")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = ""

	return cmd.Run()
}
