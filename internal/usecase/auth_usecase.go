package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-tutorial/internal/domain"
)

var (
	ErrEmailExists       = errors.New("email already registered")
	ErrInvalidCredential = errors.New("invalid email or password")
)

type AuthUsecase interface {
	Register(email, password string) (*domain.User, error)
	Login(email, password string) (string, error)
	GetUserByID(id uint) (*domain.User, error)
	ParseToken(tokenStr string) (uint, error)
}

type authUsecase struct {
	userRepo domain.UserRepository
	jwtKey   []byte
	ttl      time.Duration
}

func NewAuthUsecase(userRepo domain.UserRepository) AuthUsecase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// fallback nếu quên set env
		secret = "dev-secret-key"
	}

	return &authUsecase{
		userRepo: userRepo,
		jwtKey:   []byte(secret),
		ttl:      24 * time.Hour,
	}
}

func (u *authUsecase) Register(email, password string) (*domain.User, error) {
	// check trùng email
	if existing, _ := u.userRepo.GetByEmail(email); existing != nil {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredential
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredential
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   now.Add(u.ttl).Unix(),
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(u.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (u *authUsecase) GetUserByID(id uint) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *authUsecase) ParseToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return u.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidCredential
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidCredential
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, ErrInvalidCredential
	}

	return uint(sub), nil
}
