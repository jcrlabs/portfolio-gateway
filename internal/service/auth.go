package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"github.com/jonathanCaamano/portfolio-gateway/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users *repository.UserRepo
	cfg   *config.Config
}

func NewAuthService(users *repository.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{users, cfg}
}

type accessClaims struct {
	Email string      `json:"email"`
	Role  domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	accessToken, err = s.signAccess(user)
	if err != nil {
		return
	}
	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err = s.users.SaveRefreshToken(rt); err != nil {
		return
	}
	return accessToken, rt.Token, nil
}

func (s *AuthService) Refresh(token string) (string, error) {
	rt, err := s.users.FindRefreshToken(token)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		return "", errors.New("invalid refresh token")
	}
	user, err := s.users.FindByID(rt.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}
	return s.signAccess(user)
}

func (s *AuthService) Logout(token string) error {
	return s.users.DeleteRefreshToken(token)
}

func (s *AuthService) signAccess(user *domain.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(s.cfg.JWTSecret))
}

func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(h), err
}
