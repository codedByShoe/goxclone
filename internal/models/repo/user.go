package repo

import (
	"github.com/codedbyshoe/goxclone/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetUser(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
