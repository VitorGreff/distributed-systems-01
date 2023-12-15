package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"trab01/db"
	"trab01/models"

	"github.com/gin-gonic/gin"
)

// checked
func GetUsers(c *gin.Context) {
	var usersWithoutPassword []models.UserResponse

	for _, user := range db.Users {
		userWP := models.UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}
		usersWithoutPassword = append(usersWithoutPassword, userWP)
	}

	c.JSON(http.StatusOK, usersWithoutPassword)
}

// checked
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, user := range db.Users {
		if user.Id == id {
			userResponse := models.UserResponse{
				Id:    user.Id,
				Name:  user.Name,
				Email: user.Email,
			}
			c.JSON(http.StatusOK, userResponse)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

// checked
func DeleteUser(c *gin.Context) {
	if err := validateToken(c); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}

	for i, user := range db.Users {
		// if user.Id == id && user.Email == dto.Email && user.Password == dto.Password {
		if user.Id == id {
			db.Users = append(db.Users[:i], db.Users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Usuario com id %v deletado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Não foi possivel deletar usuario de id %v!", id)})
}

// checked
func PostUsers(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}
	newUser.Id = db.Users[len(db.Users)-1].Id + 1
	db.Users = append(db.Users, newUser)
	c.JSON(http.StatusOK, gin.H{"Resposta": "Usuário adicionado"})
}

// checked
func EditUser(c *gin.Context) {
	if err := validateToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("Erro: %v", err))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}

	var jsonData models.User
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}

	for i, u := range db.Users {
		if u.Id == id {
			updateUser(&db.Users[i], jsonData)
			c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Usuario com Id %v editado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Erro ao editar usuario de Id %v!", id)})
}

// checked
func Login(c *gin.Context) {
	var dto models.AuthDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Erro: %v", err))
		return
	}
	for _, v := range db.Users {
		if v.Email == dto.Email && v.Password == dto.Password {
			request, err := http.Get("http://localhost:8081/usuarios/token")
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro: %v", err))
				return
			}
			defer request.Body.Close()

			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro: %v", err))
				return
			}
			var token string
			err = json.Unmarshal(bodyBytes, &token)
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro: %v", err))
				return
			}

			c.JSON(http.StatusAccepted, token)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": "Dados de login inválidos"})
}

func updateUser(user *models.User, jsonData models.User) {
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

func validateToken(c *gin.Context) error {
	header := strings.Split(c.GetHeader("Authorization"), " ")
	if len(header) < 2 {
		return errors.New("token em branco")
	}

	requestToken := header[1]

	req, err := http.NewRequest("POST", "http://localhost:8081/usuarios/validar-token", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+requestToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("erro ao executar a requisição")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("requisição não autorizada")
	}
	return nil
}
