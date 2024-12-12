package common

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/patos-ufscar/quack-week/models"
)

func WhyTokenIsExpired(tokenString string, secretKey []byte) string {
	claims := &models.JwtClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// fmt.Printf("token valid: %t\n", token.Valid)

	if err != nil {
		fmt.Println("Error parsing token: ", err)
		return err.Error()
	}

	if claims.ExpiresAt > time.Now().Unix() {
		return "token is valid"
	}

	return "token timedout"
}
