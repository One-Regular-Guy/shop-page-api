package usecase

import "github.com/One-Regular-Guy/shop-page-api/internal/domain"

type RegisterUserInput struct {
	Name     string
	Username string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	Message string
}

type RegisterUserUseCase struct {
	UserRepo domain.UserRepository
}

func (uc *RegisterUserUseCase) Execute(input RegisterUserInput) (RegisterUserOutput, error) {
	usernameTaken, err := uc.UserRepo.IsUsernameTaken(input.Username)
	if err != nil {
		return RegisterUserOutput{}, err
	}
	if usernameTaken {
		return RegisterUserOutput{Message: "Username already exists"}, nil
	}

	emailTaken, err := uc.UserRepo.IsEmailTaken(input.Email)
	if err != nil {
		return RegisterUserOutput{}, err
	}
	if emailTaken {
		return RegisterUserOutput{Message: "Email already exists"}, nil
	}

	user := &domain.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password, // hash no repo infra
	}
	err = uc.UserRepo.CreateUser(user)
	if err != nil {
		return RegisterUserOutput{}, err
	}
	return RegisterUserOutput{Message: "User registered successfully"}, nil
}
