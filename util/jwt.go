package util

import (
	"avolta/config"
	"avolta/object/auth"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(email string, userId string) (string, error) {

	claims := auth.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Id:        userId,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.App.Server.ExpiredToken)).Unix(), // Token expires in 1 hour
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.App.Server.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
