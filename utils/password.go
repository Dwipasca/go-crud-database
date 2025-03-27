package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	// before we hash the password, we need to convert it to a byte slice
	// because the bcrypt.GenerateFromPassword function only accepts a byte slice
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}
func CheckPassword(hashedPassword, plainPassword string) bool {
	// CompareHashAndPassword returns nil on success and an error on failure
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
