package models

import "time"

// URL представляет структуру для хранения информации о длинных и коротких ссылках
type URL struct {
	LongLink       string
	ShortLink      string
	ExpirationTime time.Time
}
