// service/user_service.go
package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jinwoole/worklog-backend/models"
	"github.com/jinwoole/worklog-backend/repository"
	"github.com/jinwoole/worklog-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines user-related business logic
type UserService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

// userService is a concrete implementation of UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService constructs a UserService with given repository
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register creates a new user after ensuring no duplicate
func (s *userService) Register(email, password string) (*models.User, error) {
	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login authenticates user and returns JWT token
func (s *userService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
