package usecase

import (
	"github.com/One-Regular-Guy/shop-page-api/internal/domain"
)

type LoginUserInput struct {
	Username string
	Password string
}

type LoginUserOutput struct {
	Token   string
	Message string
}

type LoginUserUseCase struct {
	UserRepo      domain.UserRepository
	TokenProvider TokenProvider
}

type TokenProvider interface {
	GenerateToken(username string) (string, error)
}

func (uc *LoginUserUseCase) Execute(input LoginUserInput) (LoginUserOutput, error) {
	userExists, err := uc.UserRepo.IsUsernameTaken(input.Username)
	if err != nil {
		return LoginUserOutput{}, err
	}
	if !userExists {
		return LoginUserOutput{Message: "Credentials aren't valid"}, nil
	}
	valid, err := uc.UserRepo.CheckPassword(input.Username, input.Password)
	if err != nil {
		return LoginUserOutput{}, err
	}
	if !valid {
		return LoginUserOutput{Message: "Credentials aren't valid"}, nil
	}
	token, err := uc.TokenProvider.GenerateToken(input.Username)
	if err != nil {
		return LoginUserOutput{}, err
	}
	return LoginUserOutput{Token: token, Message: "Login successful"}, nil
}
