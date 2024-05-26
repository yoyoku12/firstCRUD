package controllers

import (
	"URL_SHORT/internal/usecases/url"
	"database/sql"
	"log"
	"net/http"
)

type URLController struct {
	URLUsecase *url.URLUsecase
}

func NewURLController(usecase *url.URLUsecase) *URLController {
	return &URLController{URLUsecase: usecase}
}

func (ctrl *URLController) LongToShort(w http.ResponseWriter, r *http.Request) {
	longLink := r.URL.Query().Get("longlink")
	if longLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Поле не может быть пустым"))
		return
	}

	shortLink, err := ctrl.URLUsecase.CreateShortURL(longLink)
	if err != nil {
		log.Println("Ошибка внесения данных в бд", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания короткой ссылки"))
		return
	}

	log.Println("127.0.0.1/" + shortLink)
	w.Write([]byte(("127.0.0.1/" + shortLink)))
}

func (ctrl *URLController) Redirect(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Path[1:]
	longLink, err := ctrl.URLUsecase.GetLongURL(shortLink)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println("Ошибка при запросе к базе данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка сервера"))
		}
		return
	}
	http.Redirect(w, r, longLink, http.StatusFound)
}
