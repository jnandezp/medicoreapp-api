package service

import (
	"errors"
	"log"

	"github.com/jnandezp/medicoreapp-api/internal/auth/model"
	"github.com/jnandezp/medicoreapp-api/internal/auth/repository"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// UserService define la interfaz para la lógica de negocio de usuarios.
// Es similar a los Services en Laravel tambien como los Traits en PHP.
type UserService interface {
	CreateUser(name, email, password string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(id uint, name string) (*model.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

var ErrEmailExists = errors.New("email already exists")

// NewUserService crea una nueva instancia del servicio de usuarios.
// Nota cómo depende de la interfaz del repositorio, no de la implementación.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser hashea la contraseña y le pide al repositorio que guarde el nuevo usuario.
func (s *userService) CreateUser(name, email, password string) (*model.User, error) {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		// Si err es nil, significa que encontró un usuario, por lo tanto, el email ya existe.
		return nil, ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Si el error es algo diferente a "no encontrado", es un error real de la base de datos.
		return nil, err
	}

	// Hasheamos la contraseña para nunca guardarla en texto plano.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Creamos la entidad de usuario.
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Le pedimos al repositorio que lo guarde.
	if err := s.repo.Create(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id uint, name string) (*model.User, error) {
	// Primero, encontramos al usuario.
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err // El usuario no existe
	}

	// Actualizamos los campos.
	user.Name = name

	// Guardamos los cambios.
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
