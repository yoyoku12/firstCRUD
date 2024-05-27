package controllers

import (
	"URL_SHORT/internal/usecases/url"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

type URLController struct {
	URLUseCase *url.Usecase
}

func NewURLController(useCase *url.Usecase) *URLController {
	return &URLController{URLUseCase: useCase}
}

func (ctrl *URLController) LongToShort(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGet(w, r)
	case http.MethodPost:
		ctrl.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Метод не поддерживается"))
		if err != nil {
			log.Println("Ошибка при отправке ответа:", err)
		}
	}
}

func (ctrl *URLController) handleGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Path[len("/"):]

	log.Println("Получен GET-запрос для короткой ссылки:", shortLink)

	longLink, err := ctrl.URLUseCase.GetLongURL(shortLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Короткая ссылка не найдена в базе данных:", shortLink)
			http.NotFound(w, r)
		} else {
			log.Println("Ошибка при запросе к базе данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Ошибка сервера"))
			if err != nil {
				log.Println("Ошибка при отправке ответа:", err)
			}
		}
		return
	}

	log.Println("Редирект на длинную ссылку:", longLink)
	http.Redirect(w, r, longLink, http.StatusFound)
}

func (ctrl *URLController) handlePost(w http.ResponseWriter, r *http.Request) {
	longLink := r.URL.Query().Get("longLink")
	if longLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Длинная ссылка не может быть пустой"))
		if err != nil {
			log.Println("Ошибка при отправке ответа:", err)
		}
		return
	}

	shortLink, err := ctrl.URLUseCase.CreateShortURL(longLink)
	if err != nil {
		log.Println("Ошибка внесения данных в бд", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Ошибка создания короткой ссылки"))
		if err != nil {
			log.Println("Ошибка при отправке ответа:", err)
		}
		return
	}

	response := "127.0.0.1/" + shortLink
	log.Println(response)
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Println("Ошибка при отправке ответа:", err)
	}
}
