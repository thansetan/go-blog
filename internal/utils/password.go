package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) ([]byte, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(raw), 5)
	if err != nil {
		return nil, err
	}
	return password, nil
}

func IsValidPassword(password []byte, raw string) error {
	return bcrypt.CompareHashAndPassword(password, []byte(raw))
}
