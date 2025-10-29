package advanced

import (
	"errors"
	"fmt"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) CreateUser(user *User) error {
	if err := s.validateUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *UserService) UpdateUser(user *User) error {
	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if err := s.validateUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	existing, err := s.repo.GetByID(user.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if existing == nil {
		return errors.New("user does not exist")
	}

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) validateUser(user *User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Age < 0 {
		return errors.New("age cannot be negative")
	}

	if user.Age > 150 {
		return errors.New("age is too high")
	}

	return nil
}
