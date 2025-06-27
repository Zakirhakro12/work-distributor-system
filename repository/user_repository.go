package repository

import (
	"work-distributor-system/models"

	"gorm.io/gorm"
)

// UserRepository provides methods for accessing User data from the database
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepo initializes a new UserRepository with a GORM database instance
func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser adds a new user to the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

// it retrieves a user by their unique username.
// Returns an error if the user is not found.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
