package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var JwtSecret = []byte("jsdjnsndnfDJNFjafkldlflfklwlkeiejfnjvnv") // ik thats need to be secure

func GenerateJWT(ID, Role string) (string, error) {
	Claims := &CustomClaims{
		Role: Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err

	}
	return tokenString, nil
}

func VerifyJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
