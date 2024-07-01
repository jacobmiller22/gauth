package service

import (
	"fmt"
	"math/rand"
)

func GenerateBearerToken(ClientId string, ClientSecret string) string {
	token := randSeq(20)
	return fmt.Sprintf("Bearer %s", token)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
