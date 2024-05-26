package main

import (
	"URL_SHORT/internal/interfaces/controllers"
	"URL_SHORT/internal/interfaces/database/postgres"
	"URL_SHORT/internal/usecases/url"
	"URL_SHORT/pkg/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка .env файла из корневой директории проекта
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("DB_USER:", os.Getenv("DB_USER")) // Добавьте эту строку для отладки

	err = config.RunMigrations()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	db, err := postgres.ConnectToDB(config.LoadDBConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	urlRepo := postgres.NewURLRepository(db)
	urlUsecase := url.NewURLUsecase(urlRepo)
	urlController := controllers.NewURLController(urlUsecase)

	go func() {
		for {
			urlUsecase.DeleteExpiredLinks()
			time.Sleep(24 * time.Hour)
		}
	}()

	http.HandleFunc("/longToShort", urlController.LongToShort)
	http.HandleFunc("/", urlController.Redirect)

	log.Println("Server starting...")
	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		log.Println("Server error", err)
	}
}
