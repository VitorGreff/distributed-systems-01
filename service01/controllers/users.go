package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	header := strings.Split(c.GetHeader("Authorization"), " ")
	if len(header) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"resposta": "Token não fornecido"})
		return
	}

	requestToken := header[1]
	resp, err := http.Post("http://localhost:8081/usuarios/validar-token", "application/json", bytes.NewBuffer([]byte(requestToken)))
	fmt.Println(resp.StatusCode)
	if err != nil || resp.StatusCode == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"resposta": "Não autorizado"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"resposta": "Id inválido"})
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

func Login(c *gin.Context) {
	// o pacote JSON precisa que os dados sejam públicos
	type UserDto struct {
		Email    string
		Password string
	}

	var dto UserDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"resposta": "Body inválido"})
	}
	for _, v := range db.Users {
		if v.Email == dto.Email && v.Password == dto.Password {
			request, err := http.Get("http://localhost:8081/usuarios/token")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"resposta": "Erro ao fazer a requisição"})
				return
			}

			defer request.Body.Close()
			var responseData map[string]interface{}
			err = json.NewDecoder(request.Body).Decode(&responseData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"resposta": "Erro ao decodificar o JSON"})
				return
			}

			token, ok := responseData["token"]
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"resposta": "Campo token não encontrado"})
				return
			}

			c.JSON(http.StatusAccepted, token)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"resposta": "Dados de login inválidos"})
}
