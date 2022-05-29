package utils

import "github.com/golang-jwt/jwt/v4"

func ParseToken(token string) (*jwt.Token, error) {
	claims, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)

	if err != nil {
		return nil, err
	}

	return claims, nil
}
