package main

import (
	"trab01/service01/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/usuarios", controllers.GetUsers)
	router.GET("/usuarios/:id", controllers.GetUser)
	router.POST("/usuarios", controllers.PostUsers)
	router.PUT("/usuarios/:id", controllers.EditUser)
	router.DELETE("/usuarios/:id", controllers.DeleteUser)

	router.Run(":5000")
}
