package main

import (
	"URL_SHORT/internal/interfaces/controllers"
	"URL_SHORT/internal/interfaces/database/postgres"
	"URL_SHORT/internal/usecases/url"
	"URL_SHORT/pkg/config"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment variables loaded")

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
	urlUseCase := url.NewURLUsecase(urlRepo)
	urlController := controllers.NewURLController(urlUseCase)

	go func() {
		t := time.NewTicker(24 * time.Hour)
		for range t.C {
			urlUseCase.DeleteExpiredLinks()
		}
	}()

	http.HandleFunc("/", urlController.LongToShort)
	http.HandleFunc("/longToShort", urlController.LongToShort)

	log.Println("Server starting...")
	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		log.Println("Server error", err)
	}
}
