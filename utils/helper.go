package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func HashPassword(password string) string {
	hashBytes := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hashBytes[:])
	return passwordHash
}

func Comparehash(password, ex_passwordhash string) bool {
	hashBytes := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hashBytes[:])
	return ex_passwordhash == passwordHash
}

func GenerateSessionToken() (*string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return nil, errors.New("Failed to create token")
	}

	tokenStr := hex.EncodeToString(token)
	return &tokenStr, nil
}
