package repositories

import (
	"github.com/jnandezp/medicoreapp-api/internal/user/entities"
	"gorm.io/gorm"
)

// UserRepository define la interfaz para las operaciones de la base de datos de usuarios.
// Es una buena práctica para permitir el testing (mocking).
// Es similar a los Models en Laravel.
type UserRepository interface {
	Create(user *entities.User) error
	FindAll() ([]entities.User, error)
	FindByID(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
}

// userRepository es la implementación concreta que usa GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos.
func (r *userRepository) Create(user *entities.User) error {
	// GORM se encarga de la consulta SQL "INSERT INTO users..."
	result := r.db.Create(user)

	return result.Error
}

func (r *userRepository) FindAll() ([]entities.User, error) {
	var users []entities.User
	result := r.db.Find(&users)

	return users, result.Error
}

func (r *userRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, id)

	return &user, result.Error
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	// GORM devolverá un error 'gorm.ErrRecordNotFound' si no encuentra al usuario.
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entities.User{}, id).Error
}
