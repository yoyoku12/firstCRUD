package controllers

import (
	"URL_SHORT/internal/domain/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthController struct {
	DB *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	getUsername := r.URL.Query().Get("login")
	getPassword := r.URL.Query().Get("password")
	getEmail := r.URL.Query().Get("email")

	if getUsername == "" || getPassword == "" || getEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Поля login, password и email не могут быть пустыми"))
		if err != nil {
			log.Println("w.write error", err)
		}
		return
	}

	user := models.User{
		Username:    getUsername,
		Email:       getEmail,
		CreatedAt:   time.Now(),
		IsActivated: false,
		Version:     1,
	}

	err := user.Password.Set(getPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Ошибка хэширования пароля"))
		if err != nil {
			log.Println("w.write error", err)
		}
		return
	}

	tx, err := ac.DB.Begin()
	if err != nil {
		log.Println("Ошибка создания транзакции:", err)
		w.Write([]byte("Ошибка создания транзакции"))
		return
	}

	_, err = tx.Exec("INSERT INTO users (username, password, email, created_at, is_activated, version) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Username, user.Password.Hash(), user.Email, user.CreatedAt, user.IsActivated, user.Version)
	if err != nil {
		tx.Rollback()
		log.Println("Ошибка создания пользователя:", err)
		w.Write([]byte(fmt.Sprintf("Ошибка создания пользователя: %v", err)))
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Ошибка коммита транзакции:", err)
		w.Write([]byte("Ошибка коммита транзакции"))
		return
	}

	log.Println("Пользователь создан")
	w.Write([]byte("Пользователь создан"))
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	getLogin := r.URL.Query().Get("login")
	getPassword := r.URL.Query().Get("password")

	if getLogin == "" || getPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Поля не могут быть пустыми"))
		if err != nil {
			log.Println("w.write error", err)
		}
		return
	}

	var user models.User
	query := "SELECT id, username, password, created_at, email, is_activated, version FROM users WHERE username = $1"
	row := ac.DB.QueryRow(query, getLogin)
	var passwordHash []byte
	err := row.Scan(&user.ID, &user.Username, &passwordHash, &user.CreatedAt, &user.Email, &user.IsActivated, &user.Version)
	if err == sql.ErrNoRows {
		w.Write([]byte("Пользователь не найден"))
		log.Println("Пользователь не найден")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	user.Password.SetHash(passwordHash)

	ok, err := user.Password.Check(getPassword)
	if !ok || err != nil {
		w.Write([]byte("Неправильный пароль"))
		log.Println("Неправильный пароль")
	} else {
		w.Write([]byte("Авторизация успешна"))
		log.Println("Авторизация успешна")
	}
}
