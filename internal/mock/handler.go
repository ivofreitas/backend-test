package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type HandlerMock struct {
	mock.Mock
}

func (m *HandlerMock) Create(ctx context.Context, param interface{}) (interface{}, error) {
	args := m.Called(ctx, param)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0), resultError
}

func (m *HandlerMock) GetByID(ctx context.Context, param interface{}) (interface{}, error) {
	args := m.Called(ctx, param)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0), resultError
}

func (m *HandlerMock) List(ctx context.Context, param interface{}) (interface{}, error) {
	args := m.Called(ctx, param)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0), resultError
}

func (m *HandlerMock) Update(ctx context.Context, param interface{}) (interface{}, error) {
	args := m.Called(ctx, param)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0), resultError
}

func (m *HandlerMock) Delete(ctx context.Context, param interface{}) (interface{}, error) {
	args := m.Called(ctx, param)

	var resultError error
	if args.Get(1) != nil {
		resultError = args.Get(1).(error)
	}

	return args.Get(0), resultError
}
