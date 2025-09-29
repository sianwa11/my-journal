package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userId int, tokenSecret string, expiresIn time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "my-journal",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: strconv.Itoa(userId),
	})

	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}


	return jwt, nil
}


func ValidateJWT(tokenString, tokenSecret string) (int, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return 0, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	if subject == "" {
		return 0, fmt.Errorf("subject claim is empty")
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		return 0, fmt.Errorf("invalid subject claim: %v", err)
	}

	return userId, nil
}