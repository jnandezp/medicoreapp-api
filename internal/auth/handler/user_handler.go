package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jnandezp/medicoreapp-api/internal/auth/service"

	"github.com/gin-gonic/gin"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios.
// Son similares a los Controllers en Laravel.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler crea una nueva instancia del handler de usuarios.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserInput define la estructura del JSON que esperamos recibir.
// Las etiquetas `binding:"required"` son para la validación automática de Gin.
// Es similar a las 'Form Requests' en Laravel.
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// CreateUser es el método que maneja la ruta POST /users.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input CreateUserInput

	// 1. Validamos el JSON de entrada.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Llamamos al servicio para crear el usuario.
	user, err := h.userService.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			// 409 Conflict es el código de estado ideal para este caso.
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}

		// Aquí podrías manejar errores más específicos, como email duplicado.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 3. Devolvemos una respuesta exitosa.
	// El código 201 Created es el estándar para creación de recursos.
	// La contraseña no se devuelve gracias a la etiqueta `json:"-"` en el modelo.
	c.JSON(http.StatusCreated, user)
}

// GetAllUsers maneja GET /users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID maneja GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Gin nos da el 'id' de la URL como un string.
	idParam := c.Param("id")
	// Lo convertimos a un número (uint).
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		// Este error comúnmente significa "no encontrado".
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUserInput define la estructura del JSON para actualizar.
type UpdateUserInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateUser maneja PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), input.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or failed to update"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser maneja DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or failed to delete"})
		return
	}
	// 204 No Content es una respuesta común para un delete exitoso.
	// O puedes devolver un 200 con un mensaje.
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
