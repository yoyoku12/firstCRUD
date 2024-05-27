package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    password  `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	Email       string    `json:"email"`
	IsActivated bool      `json:"isActivated"`
	Version     int       `json:"-"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(s string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &s
	p.hash = hash
	return nil
}

func (p *password) Check(s string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(s))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *password) Hash() []byte {
	return p.hash
}
