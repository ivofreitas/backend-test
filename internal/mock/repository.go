package mock

import (
	"backend-test/internal/domain"
	"context"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)

	var resultError error
	if args.Get(0) != nil {
		resultError = args.Get(0).(error)
	}

	return resultError
}

func (m *RepositoryMock) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	args := m.Called(ctx, id)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0).(*domain.User), resultError
}

func (m *RepositoryMock) ListUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0).([]domain.User), resultError
}

func (m *RepositoryMock) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)

	var resultError error
	if args.Get(0) != nil {
		resultError = args.Get(0).(error)
	}

	return resultError
}

func (m *RepositoryMock) DeleteUser(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)

	var resultError error
	if args.Get(0) != nil {
		resultError = args.Get(0).(error)
	}

	return resultError
}
