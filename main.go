package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPassword, dbName, dbSSLMode)

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

func generateRandStr() string {

	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"

	rand.Seed(time.Now().UnixNano())

	var result string
	for i := 0; i < 8; i++ {
		randType := rand.Intn(2)
		if randType == 0 {
			result = result + string(alphabet[rand.Intn(len(alphabet))])
		} else {
			result = result + string(numbers[rand.Intn(len(numbers))])
		}
	}

	return result
}

func deleteLinks(db *sql.DB) {
	_, err := db.Exec("DELETE FROM urls WHERE expiration_time < NOW()")
	if err != nil {
		log.Println("Ошибка при удалении старых ссылок:", err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := connectToDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go func() {
		for {
			deleteLinks(db)
			time.Sleep(24 * time.Hour)
		}
	}()

	http.HandleFunc("/longToShort", func(w http.ResponseWriter, r *http.Request) {
		getLink := r.URL.Query().Get("longlink")
		if getLink == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Поле не может быть пустым"))
			return
		}

		longLink := getLink
		shortLink := generateRandStr()
		log.Println("127.0.0.1/" + shortLink)

		expirationTime := time.Now().Add(7 * 24 * time.Hour)
		query := "INSERT INTO urls (long_link, short_link, expiration_time) VALUES ($1, $2, $3)"
		_, err = db.Exec(query, longLink, shortLink, expirationTime)
		if err != nil {
			log.Println("Ошибка внесения данных в бд", err)
			return
		}

		log.Println("Данные успешно внесены в бд")

		w.Write([]byte(("127.0.0.1/" + shortLink)))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortLink := r.URL.Path[1:]
		var longLink string
		err := db.QueryRow("SELECT long_link FROM urls WHERE short_link = $1", shortLink).Scan(&longLink)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				log.Println("Ошибка при запросе к базе данных:", err)
			}
			return
		}
		http.Redirect(w, r, longLink, http.StatusFound)
	})

	log.Println("Server starting...")
	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		log.Println("Server error", err)
	}
}
