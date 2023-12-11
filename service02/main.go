package main

import (
	"net/http"
	"trab01/service02/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.POST("/usuarios/login", func(c *gin.Context) {
		token, err := controllers.GenerateToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"resposta": "não foi possível gerar o token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})

	})

	router.Run(":6000")
}
