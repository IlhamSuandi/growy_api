package auth

import "golang.org/x/crypto/bcrypt"

func ComparePassword(password string, hashedPassword string) bool {
	// Compare the password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
