package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jnandezp/medicoreapp-api/internal/user/entities"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	var dsn string

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error: No se pudo cargar el archivo .env. Asegúrate de que exista en la raíz del proyecto.")
	}

	// 2. Lee la variable para decidir qué driver usar
	dbConnection := os.Getenv("DB_CONNECTION")
	log.Println(dbConnection)
	log.Printf("Usando la conexión de base de datos: %s", dbConnection)

	// 3. Usa un 'switch' para manejar la lógica de conexión
	switch dbConnection {
	case "mysql":
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbname := os.Getenv("DB_NAME")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbname,
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "sqlite":
		dbName := os.Getenv("DB_DATABASE")
		DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	default:
		log.Fatalf("Conexión de base de datos no válida: %s", dbConnection)
	}

	if err != nil {
		log.Fatalf("Falló la conexión a la base de datos '%s': %v", dbConnection, err)
	}

	// El resto del código es el mismo para ambas bases de datos
	err = DB.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatal("Falló la migración:", err)
	}

	log.Println("Conexión y migración exitosas.")
}
