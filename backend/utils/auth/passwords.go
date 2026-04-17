package authutils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func VerificatePassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
