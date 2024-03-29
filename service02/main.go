package main

import (
	"net/http"
	"trab01/service02/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/usuarios/token", func(c *gin.Context) {
		token, err := controllers.GenerateToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Resposta": "Não foi possível gerar o token"})
			return
		}
		c.String(http.StatusOK, token)
	})

	router.POST("/usuarios/validar-token", func(c *gin.Context) {
		if err := controllers.ValidateToken(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Resposta": "Token inválido"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"Resposta": "Token válido"})
	})

	router.Run(":8081")
}
