package url

import (
	"URL_SHORT/internal/domain/repositories"
	"URL_SHORT/pkg/utils"
	"time"
)

type URLUsecase struct {
	URLRepo repositories.URLRepository
}

func NewURLUsecase(repo repositories.URLRepository) *URLUsecase {
	return &URLUsecase{URLRepo: repo}
}

func (uc *URLUsecase) CreateShortURL(longLink string) (string, error) {
	shortLink := utils.GenerateRandStr()
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	err := uc.URLRepo.SaveURL(longLink, shortLink, expirationTime)
	return shortLink, err
}

func (uc *URLUsecase) GetLongURL(shortLink string) (string, error) {
	return uc.URLRepo.GetLongURL(shortLink)
}

func (uc *URLUsecase) DeleteExpiredLinks() {
	uc.URLRepo.DeleteExpiredLinks()
}
