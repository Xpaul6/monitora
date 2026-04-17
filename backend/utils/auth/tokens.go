package authutils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
