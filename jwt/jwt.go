package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("jsdjnsndnfDJNFjafkldlflfklwlkeiejfnjvnv")

func GenerateJWT(ID string, Role string) (string, error) {
	claims := &CustomClaims{
		Role: Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err

	}
	return tokenString, nil
}
