package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"trab01/db"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, db.Users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, user := range db.Users {
		if user.Id == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for i, user := range db.Users {
		if user.Id == id {
			db.Users = append(db.Users[:i], db.Users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"resposta": fmt.Sprintf("Usuario com Id %v deletado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

func PostUsers(c *gin.Context) {
	var newUser db.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	db.Users = append(db.Users, newUser)
	c.JSON(http.StatusOK, db.Users)
}

func EditUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var jsonData db.User
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for i, u := range db.Users {
		if u.Id == id {
			updateUser(&db.Users[i], jsonData)
			c.JSON(http.StatusOK, gin.H{"resposta": fmt.Sprintf("Usuario com Id %v editado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

func updateUser(user *db.User, jsonData db.User) {
	if jsonData.Name != "" {
		user.Name = jsonData.Name
	}
	if jsonData.Email != "" {
		user.Email = jsonData.Email
	}
	if jsonData.Password != "" {
		user.Password = jsonData.Password
	}
}

// func Login(c *gin.Context) {
// 	type passwordChange struct {
// 		email        string
// 		old_password string
// 		new_password string
// 	}

// 	var pc passwordChange

// 	if err := c.ShouldBindJSON(&pc); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"resposta": "Body inválido"})
// 	}

// 	for _, v := range db.Users {
// 		if v.Email == pc.email {

// 		}
// 	}
// }
