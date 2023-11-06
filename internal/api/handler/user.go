package handler

import (
	"backend-test/internal/domain"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type GenerateFromPassword func(password []byte, cost int) ([]byte, error)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id int32) error
	GetUser(ctx context.Context, id int32) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
}

type User struct {
	repository           Repository
	generateFromPassword GenerateFromPassword
}

// NewUserHdl - create a new instance of user handler
func NewUserHdl(repository Repository, generateFromPassword GenerateFromPassword) *User {
	return &User{repository, generateFromPassword}
}

// Create
// @Summary create a user.
// @Param key body domain.User true "request body"
// @Tags User
// @Accept json
// @Product json
// @Success 201 {object} domain.UserNoPass
// @Failure 400 {object} domain.ResponseError
// @Failure 500 {object} domain.ResponseError
// @Router /user [post]
func (u *User) Create(ctx context.Context, param interface{}) (interface{}, error) {
	user := param.(*domain.User)

	hash, err := u.generateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrorDiscover(domain.BadRequest{DeveloperMessage: err.Error()})
	}
	user.Password = string(hash)

	err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	return domain.NewUserWithoutPass(user), nil
}

// GetByID
// @Summary retrieve a user.
// @Param id path string true "User id"
// @Tags User
// @Accept json
// @Product json
// @Success 200 {object} domain.UserNoPass
// @Failure 400 {object} domain.ResponseError
// @Failure 404 {object} domain.ResponseError
// @Failure 500 {object} domain.ResponseError
// @Router /user/{id} [get]
func (u *User) GetByID(ctx context.Context, param interface{}) (interface{}, error) {
	request := param.(*domain.GetByIDRequest)
	user, err := u.repository.GetUser(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrorDiscover(domain.NotFound{DeveloperMessage: err.Error()})
		}

		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	return domain.NewUserWithoutPass(user), nil
}

// List
// @Summary list all users.
// @Tags User
// @Product json
// @Success 200 {object} []domain.UserNoPass
// @Failure 400 {object} domain.ResponseError
// @Failure 404 {object} domain.ResponseError
// @Failure 500 {object} domain.ResponseError
// @Router /user [get]
func (u *User) List(ctx context.Context, param interface{}) (interface{}, error) {
	users, err := u.repository.ListUsers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrorDiscover(domain.NotFound{DeveloperMessage: err.Error()})
		}

		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	response := make([]domain.UserNoPass, 0, len(users))
	for _, user := range users {
		response = append(response, *domain.NewUserWithoutPass(&user))
	}

	return response, nil
}

// Update
// @Summary update a user.
// @Param id path string true "User id"
// @Param key body domain.UpdateRequest true "request body"
// @Tags User
// @Product json
// @Success 200 {object} domain.UserNoPass
// @Failure 400 {object} domain.ResponseError
// @Failure 404 {object} domain.ResponseError
// @Failure 500 {object} domain.ResponseError
// @Router /user/{id} [put]
func (u *User) Update(ctx context.Context, param interface{}) (interface{}, error) {
	updateRequest := param.(*domain.UpdateRequest)

	user := &domain.User{
		ID:       updateRequest.ID,
		Name:     updateRequest.Name,
		Age:      updateRequest.Age,
		Email:    updateRequest.Email,
		Password: updateRequest.Password,
		Address:  updateRequest.Address,
	}

	err := u.repository.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrorDiscover(domain.NotFound{DeveloperMessage: err.Error()})
		}

		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	result, err := u.repository.GetUser(ctx, updateRequest.ID)
	if err != nil {
		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	return domain.NewUserWithoutPass(result), nil
}

// Delete
// @Summary delete a user.
// @Param id path string true "User id"
// @Tags User
// @Product json
// @Success 204
// @Failure 400 {object} domain.ResponseError
// @Failure 404 {object} domain.ResponseError
// @Failure 500 {object} domain.ResponseError
// @Router /user/{id} [delete]
func (u *User) Delete(ctx context.Context, param interface{}) (interface{}, error) {
	deleteRequest := param.(*domain.DeleteRequest)

	_, err := u.repository.GetUser(ctx, deleteRequest.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrorDiscover(domain.NotFound{DeveloperMessage: err.Error()})
		}

		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	err = u.repository.DeleteUser(ctx, deleteRequest.ID)
	if err != nil {
		return nil, domain.ErrorDiscover(domain.ResponseError{DeveloperMessage: err.Error()})
	}

	return nil, nil
}
