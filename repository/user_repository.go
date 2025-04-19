package repository

import (
	"github.com/jinwoole/worklog-backend/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository returns a UserRepository backed by sqlx.DB
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	// Insert user and return generated ID and timestamp
	return r.db.QueryRowx(
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`,
		user.Email, user.PasswordHash,
	).StructScan(user)
}
