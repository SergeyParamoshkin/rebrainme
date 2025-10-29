package advanced

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		expectedUser := &User{ID: 1, Name: "John", Email: "john@example.com", Age: 30}

		mockRepo.On("GetByID", 1).Return(expectedUser, nil)

		service := NewUserService(mockRepo)
		user, err := service.GetUser(1)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		user, err := service.GetUser(0)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", 1).Return(nil, errors.New("database error"))

		service := NewUserService(mockRepo)
		user, err := service.GetUser(1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		setupMock   func(*MockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			user: &User{Name: "Alice", Email: "alice@example.com", Age: 25},
			setupMock: func(m *MockUserRepository) {
				m.On("Create", mock.AnythingOfType("*advanced.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "missing name",
			user:        &User{Email: "test@example.com", Age: 25},
			setupMock:   func(m *MockUserRepository) {},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:        "missing email",
			user:        &User{Name: "Bob", Age: 25},
			setupMock:   func(m *MockUserRepository) {},
			wantErr:     true,
			errContains: "email is required",
		},
		{
			name:        "negative age",
			user:        &User{Name: "Charlie", Email: "charlie@example.com", Age: -5},
			setupMock:   func(m *MockUserRepository) {},
			wantErr:     true,
			errContains: "age cannot be negative",
		},
		{
			name:        "age too high",
			user:        &User{Name: "Diana", Email: "diana@example.com", Age: 200},
			setupMock:   func(m *MockUserRepository) {},
			wantErr:     true,
			errContains: "age is too high",
		},
		{
			name: "repository error",
			user: &User{Name: "Eve", Email: "eve@example.com", Age: 30},
			setupMock: func(m *MockUserRepository) {
				m.On("Create", mock.AnythingOfType("*advanced.User")).
					Return(errors.New("database error"))
			},
			wantErr:     true,
			errContains: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			service := NewUserService(mockRepo)
			err := service.CreateUser(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		user := &User{ID: 1, Name: "John Updated", Email: "john@example.com", Age: 31}
		existingUser := &User{ID: 1, Name: "John", Email: "john@example.com", Age: 30}

		mockRepo.On("GetByID", 1).Return(existingUser, nil)
		mockRepo.On("Update", user).Return(nil)

		service := NewUserService(mockRepo)
		err := service.UpdateUser(user)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		user := &User{ID: 0, Name: "John", Email: "john@example.com", Age: 30}
		err := service.UpdateUser(user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", 999).Return(nil, errors.New("not found"))

		service := NewUserService(mockRepo)
		user := &User{ID: 999, Name: "Ghost", Email: "ghost@example.com", Age: 25}
		err := service.UpdateUser(user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		user := &User{ID: 1, Name: "", Email: "test@example.com", Age: 25}
		err := service.UpdateUser(user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Delete", 1).Return(nil)

		service := NewUserService(mockRepo)
		err := service.DeleteUser(1)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		err := service.DeleteUser(0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Delete", 1).Return(errors.New("database error"))

		service := NewUserService(mockRepo)
		err := service.DeleteUser(1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete user")
		mockRepo.AssertExpectations(t)
	})
}

func TestValidateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    &User{Name: "John", Email: "john@example.com", Age: 30},
			wantErr: false,
		},
		{
			name:    "empty name",
			user:    &User{Name: "", Email: "test@example.com", Age: 30},
			wantErr: true,
		},
		{
			name:    "empty email",
			user:    &User{Name: "John", Email: "", Age: 30},
			wantErr: true,
		},
		{
			name:    "negative age",
			user:    &User{Name: "John", Email: "john@example.com", Age: -1},
			wantErr: true,
		},
		{
			name:    "age too high",
			user:    &User{Name: "John", Email: "john@example.com", Age: 151},
			wantErr: true,
		},
		{
			name:    "zero age valid",
			user:    &User{Name: "Baby", Email: "baby@example.com", Age: 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateUser(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &User{Name: "Integration Test", Email: "test@example.com", Age: 25}

	mockRepo.On("Create", user).Return(nil).Once()

	err := service.CreateUser(user)
	require.NoError(t, err)

	user.ID = 1
	mockRepo.On("GetByID", 1).Return(user, nil).Once()

	retrieved, err := service.GetUser(1)
	require.NoError(t, err)
	assert.Equal(t, user.Name, retrieved.Name)

	user.Name = "Updated Name"
	mockRepo.On("GetByID", 1).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err = service.UpdateUser(user)
	require.NoError(t, err)

	mockRepo.On("Delete", 1).Return(nil).Once()

	err = service.DeleteUser(1)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
