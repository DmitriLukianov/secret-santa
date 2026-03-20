package services

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user domain.User) error {

	if user.Name == "" {
		return fmt.Errorf("name is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, name, email *string) error {
	return s.userRepo.UpdateUser(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
