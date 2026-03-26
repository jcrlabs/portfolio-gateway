package repository

import (
	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	return &u, r.db.Where("email = ?", email).First(&u).Error
}

func (r *UserRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	var u domain.User
	return &u, r.db.First(&u, "id = ?", id).Error
}

func (r *UserRepo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) SaveRefreshToken(rt *domain.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *UserRepo) FindRefreshToken(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	return &rt, r.db.Where("token = ?", token).First(&rt).Error
}

func (r *UserRepo) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}
