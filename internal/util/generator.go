package util

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var charsets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rng.Intn(len(charsets))]
	}

	return string(letters)
}

func GenerateRandomNumber(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var charsets = []rune("0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rng.Intn(len(charsets))]
	}

	return string(letters)
}
