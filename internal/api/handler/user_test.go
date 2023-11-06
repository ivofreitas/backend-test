package handler

import (
	"backend-test/internal/domain"
	"backend-test/internal/mock"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"testing"
)

var (
	repoErr         = errors.New("repo fail")
	expectedRepoErr = domain.ResponseError{DeveloperMessage: repoErr.Error()}
	passHashErr     = errors.New("passHash failed")
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		Name          string
		Request       *domain.User
		CreateRepoErr error
		PassHashErr   error
		ExpectedError error
	}{
		{
			Name: "Test Case 1",
			Request: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
		},
		{
			Name: "Test Case 2",
			Request: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			CreateRepoErr: repoErr,
			ExpectedError: expectedRepoErr,
		},
		{
			Name: "Test Case 3",
			Request: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			PassHashErr:   passHashErr,
			ExpectedError: passHashErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			repositoryMock := &mock.RepositoryMock{}

			repositoryMock.
				On("CreateUser", ctx, tc.Request).
				Return(tc.CreateRepoErr)

			passHashMock := mock.PassHashMock{}
			passHashMock.
				On("GenerateFromPassword", testifymock.Anything, testifymock.Anything).
				Return([]byte{}, tc.PassHashErr)

			hdl := NewUserHdl(repositoryMock, passHashMock.GenerateFromPassword)
			response, err := hdl.Create(ctx, tc.Request)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestGetByID(t *testing.T) {
	testCases := []struct {
		Name          string
		Request       *domain.GetByIDRequest
		GetRepoResult *domain.User
		GetRepoErr    error
		ExpectedError error
	}{
		{
			Name: "Test Case 1",
			Request: &domain.GetByIDRequest{
				ID: 1,
			},
			GetRepoResult: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
		},
		{
			Name: "Test Case 2",
			Request: &domain.GetByIDRequest{
				ID: 1,
			},
			GetRepoErr:    sql.ErrNoRows,
			ExpectedError: errors.New("sql: no rows in result set"),
		},
		{
			Name: "Test Case 3",
			Request: &domain.GetByIDRequest{
				ID: 1,
			},
			GetRepoErr:    errors.New("repo error"),
			ExpectedError: domain.ResponseError{DeveloperMessage: "repo error"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			repositoryMock := &mock.RepositoryMock{}

			repositoryMock.
				On("GetUser", ctx, tc.Request.ID).
				Return(tc.GetRepoResult, tc.GetRepoErr)

			passHashMock := mock.PassHashMock{}

			hdl := NewUserHdl(repositoryMock, passHashMock.GenerateFromPassword)
			response, err := hdl.GetByID(ctx, tc.Request)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		Name          string
		GetRepoResult []domain.User
		ListRepoErr   error
		ExpectedError error
	}{
		{
			Name: "Test Case 1",
			GetRepoResult: []domain.User{{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			}},
		},
		{
			Name:          "Test Case 2",
			ListRepoErr:   sql.ErrNoRows,
			ExpectedError: errors.New("sql: no rows in result set"),
		},
		{
			Name:          "Test Case 3",
			ListRepoErr:   errors.New("repo error"),
			ExpectedError: domain.ResponseError{DeveloperMessage: "repo error"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			repositoryMock := &mock.RepositoryMock{}

			repositoryMock.
				On("ListUsers", ctx).
				Return(tc.GetRepoResult, tc.ListRepoErr)

			passHashMock := mock.PassHashMock{}

			hdl := NewUserHdl(repositoryMock, passHashMock.GenerateFromPassword)
			response, err := hdl.List(ctx, nil)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		Name          string
		Request       *domain.UpdateRequest
		PassHashErr   error
		UpdateRepoErr error
		GetResponse   *domain.User
		GetRepoErr    error
		ExpectedError error
	}{
		{
			Name: "Test Case 1",
			Request: &domain.UpdateRequest{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			GetResponse: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
		},
		{
			Name: "Test Case 2",
			Request: &domain.UpdateRequest{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			UpdateRepoErr: repoErr,
			ExpectedError: expectedRepoErr,
		},
		{
			Name: "Test Case 3",
			Request: &domain.UpdateRequest{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			UpdateRepoErr: sql.ErrNoRows,
			ExpectedError: errors.New("sql: no rows in result set"),
		},
		{
			Name: "Test Case 4",
			Request: &domain.UpdateRequest{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			GetRepoErr:    repoErr,
			ExpectedError: expectedRepoErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			repositoryMock := &mock.RepositoryMock{}
			repositoryMock.
				On("UpdateUser", ctx, testifymock.Anything).
				Return(tc.UpdateRepoErr)

			repositoryMock.
				On("GetUser", ctx, tc.Request.ID).
				Return(tc.GetResponse, tc.GetRepoErr)

			passHashMock := mock.PassHashMock{}

			hdl := NewUserHdl(repositoryMock, passHashMock.GenerateFromPassword)
			response, err := hdl.Update(ctx, tc.Request)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		Name          string
		Request       *domain.DeleteRequest
		GetRepoResult *domain.User
		GetRepoErr    error
		DeleteRepoErr error
		ExpectedError error
	}{
		{
			Name: "Test Case 1",
			Request: &domain.DeleteRequest{
				ID: 1,
			},
			GetRepoResult: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
		},
		{
			Name: "Test Case 2",
			Request: &domain.DeleteRequest{
				ID: 1,
			},
			GetRepoErr:    sql.ErrNoRows,
			ExpectedError: errors.New("sql: no rows in result set"),
		},
		{
			Name: "Test Case 3",
			Request: &domain.DeleteRequest{
				ID: 1,
			},
			GetRepoErr:    errors.New("repo error"),
			ExpectedError: domain.ResponseError{DeveloperMessage: "repo error"},
		},
		{
			Name: "Test Case 4",
			Request: &domain.DeleteRequest{
				ID: 1,
			},
			DeleteRepoErr: errors.New("repo error"),
			ExpectedError: domain.ResponseError{DeveloperMessage: "repo error"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			repositoryMock := &mock.RepositoryMock{}

			repositoryMock.
				On("GetUser", ctx, tc.Request.ID).
				Return(tc.GetRepoResult, tc.GetRepoErr)

			repositoryMock.
				On("DeleteUser", ctx, tc.Request.ID).
				Return(tc.DeleteRepoErr)

			passHashMock := mock.PassHashMock{}

			hdl := NewUserHdl(repositoryMock, passHashMock.GenerateFromPassword)
			response, err := hdl.Delete(ctx, tc.Request)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, response)
		})
	}
}
