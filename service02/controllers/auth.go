package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("chaveAbsurdamenteSecreta")

func GenerateToken(c *gin.Context) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	// token expires in 6 hours
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("erro ao assinar o token")
	}

	return signedToken, nil
}

func ValidateToken(c *gin.Context) error {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, returnKey)
	if err != nil {
		return errors.New("erro ao dar parse no token")
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")
}

func extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")

	// bearer <token>
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func returnKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("metodo de assinatura inesperado! %v", token.Header["alg"])
	}
	return jwtSecret, nil
}
