package main

import (
	"URL_SHORT/internal/interfaces/controllers"
	"URL_SHORT/internal/interfaces/database/postgres"
	"URL_SHORT/internal/interfaces/middleware"
	"URL_SHORT/internal/interfaces/routes"
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

	dbConfig := config.LoadDBConfig()

	db, err := postgres.ConnectToDB(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	urlRepo := postgres.NewURLRepository(db)
	urlUseCase := url.NewURLUseCase(urlRepo)
	urlController := controllers.NewURLController(urlUseCase)
	authController := controllers.NewAuthController(db)
	authMiddleware := middleware.NewAuthMiddleware(urlRepo)

	go func() {
		t := time.NewTicker(24 * time.Hour)
		for range t.C {
			urlUseCase.DeleteExpiredLinks()
		}
	}()

	routes.RegisterRoutes(urlController, authController, authMiddleware)

	log.Println("Server starting...")
	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		log.Println("Server error", err)
	}
}
