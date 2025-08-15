package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(pasword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}

	return nil
}
