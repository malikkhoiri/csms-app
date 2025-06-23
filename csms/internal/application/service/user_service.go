package service

import (
	"context"

	"github.com/malikkhoiri/csms/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id uint, status string) error {
	return s.userRepo.UpdateStatus(ctx, id, status)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
