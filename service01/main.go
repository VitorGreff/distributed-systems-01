package main

import (
	"trab01/service01/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/usuarios", controllers.GetUsers)
	router.GET("/usuarios/:id", controllers.GetUser)
	router.POST("/usuarios", controllers.PostUser)
	router.GET("/usuarios/login", controllers.Login)
	router.PUT("/usuarios/:id", controllers.EditUser)
	router.DELETE("/usuarios/:id", controllers.DeleteUser)

	router.Run(":8080")
}
