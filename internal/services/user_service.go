package services

import (
	"context"
	"errors"

	"expense-tracker/internal/models"
	"expense-tracker/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) RegisterUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*models.User, error) {
	_, err := service.userRepository.GetUserByEmail(ctx, email)

	if err == nil {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	err = service.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (service *UserService) LoginUser(ctx context.Context, email string, password string) (*models.User, error) {

	user, err := service.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, errors.New("Invalid Email or Password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
