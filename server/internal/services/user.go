package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Farzam-Haji/Ecommerce/server/internal/auth"
	"github.com/Farzam-Haji/Ecommerce/server/internal/models"
	"github.com/Farzam-Haji/Ecommerce/server/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create (context.Context, *models.User) (error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetByID(context.Context, int64) (*models.User, error)
}

type UserService struct {
	repository UserRepository
	authentication	*auth.JWTManager
}

func NewUserService(repoInterface UserRepository, authManager *auth.JWTManager) *UserService {
	return &UserService{
		repository: repoInterface,
		authentication: authManager,
	}
}

func (us *UserService) Register(ctx context.Context, email string, password string) (*models.User, string, error) {
	_, err := us.repository.GetByEmail(ctx, email)
	switch{
	case err == nil:
		return nil, "", ErrEmailUnavailable
	case errors.Is(err, repositories.ErrRowNotFound):
		//email available
	default:
		return nil, "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if len(password) < 8 {
		return nil, "", ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	newUser := models.User{
		Email: email,
		PasswordHash: string(hash),
	}

	err = us.repository.Create(ctx, &newUser)
	if errors.Is(err, repositories.ErrDuplicateKey) {
		return nil, "", ErrEmailUnavailable
	}
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	token, err := us.authentication.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrTokenGenFailed, err)
	}
	return &newUser, token, nil	
}

func (us *UserService) Login(ctx context.Context, email string, password string) (string, error){
	user, err := us.repository.GetByEmail(ctx, email)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return "", ErrWrongCredential
	}
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrWrongCredential
	}
	
	token, err := us.authentication.GenerateToken(user.ID, email)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenGenFailed, err)
	}
	return token, nil
}