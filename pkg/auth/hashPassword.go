package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	// Hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPass, nil
}
