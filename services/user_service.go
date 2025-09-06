package services

import (
	"errors"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"gorm.io/gorm"
)

// UserService interface defines user business logic
type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	AuthenticateUser(username, password string) (*models.User, error)
}

// userService implements UserService interface
type userService struct {
	db *gorm.DB
}

// NewUserService creates a new user service instance
func NewUserService() UserService {
	return &userService{
		db: db.DB,
	}
}

// CreateUser creates a new user with hashed password
func (s *userService) CreateUser(user *models.User) error {
	// Check if username already exists
	var existingUser models.User
	if err := s.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if user.Email != "" {
		if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			return errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	return s.db.Create(user).Error
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user *models.User) error {
	// Check if user exists
	var existing models.User
	if err := s.db.First(&existing, user.ID).Error; err != nil {
		return err
	}
	return s.db.Save(user).Error
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.User{}, id).Error
}

// AuthenticateUser validates user credentials
func (s *userService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
