package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandStr() string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	chars := alphabet + numbers

	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < 8; i++ {
		randIndex := rand.Intn(len(chars))
		sb.WriteByte(chars[randIndex])
	}

	return sb.String()
}
