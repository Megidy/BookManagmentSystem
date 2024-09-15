package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func UnHashPassword(HashedPassword []byte, Password []byte) error {
	err := bcrypt.CompareHashAndPassword(HashedPassword, Password)
	if err != nil {
		return err
	}
	return nil
}
