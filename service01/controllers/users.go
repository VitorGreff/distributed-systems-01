package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"trab01/db"

	"github.com/gin-gonic/gin"
)

// o pacote JSON precisa que os dados sejam públicos
type UserDto struct {
	Email    string
	Password string
}

// checked
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, db.Users)
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
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

// checked
func DeleteUser(c *gin.Context) {
	if err := validateToken(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Erro de validação do token"})
		return
	}

	var dto UserDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Body inválido"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Id inválido"})
		return
	}

	for i, user := range db.Users {
		if user.Id == id && user.Email == dto.Email && user.Password == dto.Password {
			db.Users = append(db.Users[:i], db.Users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Usuario com Id %v deletado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

// checked
func PostUsers(c *gin.Context) {
	var newUser db.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Body inválido"})
		return
	}
	newUser.Id = db.Users[len(db.Users)-1].Id + 1
	db.Users = append(db.Users, newUser)
	c.JSON(http.StatusOK, gin.H{"Resposta": "Usuário adicionado"})
}

// checked
func EditUser(c *gin.Context) {
	if err := validateToken(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Erro de validação do token"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Id inválido"})
		return
	}

	var jsonData db.User
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Body inválido"})
		return
	}

	for i, u := range db.Users {
		if u.Id == id {
			updateUser(&db.Users[i], jsonData)
			c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Usuario com Id %v editado!", id)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": fmt.Sprintf("Id %v não está cadastrado!", id)})
}

// checked
func Login(c *gin.Context) {
	// email and password
	var dto UserDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Resposta": "Body inválido"})
	}
	for _, v := range db.Users {
		if v.Email == dto.Email && v.Password == dto.Password {
			request, err := http.Get("http://localhost:8081/usuarios/token")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Resposta": "Erro ao fazer a requisição"})
				return
			}

			defer request.Body.Close()
			var responseData map[string]interface{}
			err = json.NewDecoder(request.Body).Decode(&responseData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Resposta": "Erro ao decodificar o JSON"})
				return
			}

			token, ok := responseData["token"]
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"Resposta": "Campo token não encontrado"})
				return
			}

			c.JSON(http.StatusAccepted, token)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Resposta": "Dados de login inválidos"})
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
