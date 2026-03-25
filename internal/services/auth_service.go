package services

import (
	"errors"
	"fmt"
	"time"

	"adarel-api/internal/models"
	"adarel-api/internal/repositories"
	"adarel-api/pkg/sanitize"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(name, email, password, tenantName string) error
	Login(email, password string) (string, error)
	ParseToken(tokenString string) (*AuthClaims, error)
}

type AuthClaims struct {
	UserID   uint `json:"user_id"`
	TenantID uint `json:"tenant_id"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
	jwtSecret  string
}

func NewAuthService(userRepo repositories.UserRepository, tenantRepo repositories.TenantRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, tenantRepo: tenantRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(name, email, password, tenantName string) error {
	name = sanitize.Text(name)
	email = sanitize.Text(email)
	tenantName = sanitize.Text(tenantName)

	if name == "" || email == "" || password == "" || tenantName == "" {
		return errors.New("invalid input")
	}

	if len(password) < 8 {
		return errors.New("password must have at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash failed: %w", err)
	}

	tenant := &models.Tenant{Name: tenantName}
	if err := s.tenantRepo.Create(tenant); err != nil {
		return err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		TenantID:     tenant.ID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(email, password string) (string, error) {
	email = sanitize.Text(email)
	if email == "" || password == "" {
		return "", errors.New("invalid credentials")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := AuthClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (s *authService) ParseToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
