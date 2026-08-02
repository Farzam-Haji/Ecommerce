package services

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/models"

	// "github.com/btalv/Ecommerce/server/internal/repositories"
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
		return nil, "", ErrEmailExists
	case errors.Is(err, sql.ErrNoRows):
		//email available
	default:
		return nil, "", err
	}

	if len(password) < 8 {
		return nil, "", ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	newUser := models.User{
		Email: email,
		PasswordHash: string(hash),
	}

	err = us.repository.Create(ctx, &newUser)
	if err != nil {
		return nil, "", err
	}

	token, err := us.authentication.GenerateToken(newUser.ID)
	if err != nil {
		log.Println(err)
	}
	return &newUser, token, nil	
}

func (us *UserService) Login(ctx context.Context, email string, password string) (string, error){
	user, err := us.repository.GetByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrWrongCredential
	}
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrWrongCredential
	}
	
	return us.authentication.GenerateToken(user.ID)
}