package url

import (
	"URL_SHORT/internal/domain/repositories"
	"URL_SHORT/pkg/utils"
	"time"
)

type Usecase struct {
	URLRepo repositories.URLRepository
}

func NewURLUsecase(repo repositories.URLRepository) *Usecase {
	return &Usecase{URLRepo: repo}
}

func (uc *Usecase) CreateShortURL(longLink string) (string, error) {
	shortLink := utils.GenerateRandStr()
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	err := uc.URLRepo.SaveURL(longLink, shortLink, expirationTime)
	return shortLink, err
}

func (uc *Usecase) GetLongURL(shortLink string) (string, error) {
	return uc.URLRepo.GetLongURL(shortLink)
}

func (uc *Usecase) DeleteExpiredLinks() {
	uc.URLRepo.DeleteExpiredLinks()
}
