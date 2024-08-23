package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(username, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

type TokenClaims struct {
	UserId string `json:"userId"`
}

func getParsedToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signin method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, errors.New("could not parse token")
	}
	return parsedToken, nil
}

func isParsedTokenValid(parsedToken *jwt.Token) bool {
	isValid := parsedToken.Valid
	return isValid
}

func VerifyAccessToken(token string) (*TokenClaims, error) {
	parsedToken, err := getParsedToken(token)
	if err != nil {
		return nil, errors.New("could not parse token")
	}

	isValid := isParsedTokenValid(parsedToken)
	if !isValid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userId := claims["userId"].(string)
	return &TokenClaims{UserId: userId}, nil
}
