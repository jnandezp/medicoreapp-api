package main

import (
	"fmt"

	"github.com/jnandezp/medicoreapp-api/internal/auth/handler"
	"github.com/jnandezp/medicoreapp-api/internal/auth/repository"
	"github.com/jnandezp/medicoreapp-api/internal/auth/service"
	"github.com/jnandezp/medicoreapp-api/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	// --- Inyección de Dependencias ---
	// 1. Creamos una instancia del Repositorio, pasándole la conexión a la BD.
	userRepository := repository.NewUserRepository(database.DB)
	// 2. Creamos una instancia del Servicio, pasándole el Repositorio.
	userService := service.NewUserService(userRepository)
	// 3. Creamos una instancia del Handler, pasándole el Servicio.
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	// --- Agrupamos las rutas de usuarios ---
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)       // POST /users
		userRoutes.GET("", userHandler.GetAllUsers)       // GET /users
		userRoutes.GET("/:id", userHandler.GetUserByID)   // GET /users/123
		userRoutes.PUT("/:id", userHandler.UpdateUser)    // PUT /users/123
		userRoutes.DELETE("/:id", userHandler.DeleteUser) // DELETE /users/123
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8080")

	fmt.Println("medicoreappapi package")

}
