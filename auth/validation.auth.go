package auth

// Verification et lecture du contenu des jetons JWT.

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Verifie la signature du jeton et renvoie les informations qu'il contient.
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("methode de signature inattendue : %v", token.Header["alg"])
			}
			return jwtSecret, nil
		},
		jwt.WithIssuer("cinetalk"),
		jwt.WithAudience("cinetalk-web"),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalide")
	}

	return claims, nil
}
