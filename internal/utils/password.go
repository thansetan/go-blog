package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(raw), 5)
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func IsValidPassword(password, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(raw))
}
