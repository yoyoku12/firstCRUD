package models

import "time"

type URL struct {
	LongLink       string
	ShortLink      string
	ExpirationTime time.Time
}
