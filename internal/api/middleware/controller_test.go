package middleware

import (
	"backend-test/internal/adapter/log"
	"backend-test/internal/domain"
	"backend-test/internal/mock"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	bindErrParam = `{"name": 1,"age": 18,"email": "john.doe@gmail.com","password": "password","address": "200 Random St"}`
)

func TestHandleCreateUser(t *testing.T) {
	testCases := []struct {
		Name           string
		Param          interface{}
		HandlerErr     error
		ExpectedError  error
		ExpectedStatus int
	}{
		{
			Name: "Test Case 1",
			Param: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "Test Case 2",
			Param:          bindErrParam,
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Test Case 3",
			Param: &domain.User{
				ID:       123,
				Name:     "John",
				Age:      18,
				Password: "password",
				Address:  "20 random St",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			b, _ := json.Marshal(&tc.Param)
			req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(string(b)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e := echo.New()
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			ctx := log.InitParams(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))

			e.Binder = NewBinder()
			e.Validator = NewValidator()

			handlerMock := &mock.HandlerMock{}
			handlerMock.
				On("Create", ctx, testifymock.Anything).
				Return(tc.Param, tc.HandlerErr)

			ctrl := NewController(handlerMock.Create, http.StatusCreated, new(domain.User))
			err := ctrl.Handle(c)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.Equal(t, tc.ExpectedStatus, rec.Code)
			assert.NotNil(t, rec.Body)
			assert.NoError(t, err)
		})
	}
}

func TestHandleUpdateUser(t *testing.T) {
	testCases := []struct {
		Name           string
		Param          interface{}
		HandlerErr     error
		ExpectedError  error
		ExpectedStatus int
	}{
		{
			Name: "Test Case 1",
			Param: &domain.UpdateRequest{
				ID:       123,
				Name:     "John",
				Age:      18,
				Email:    "john@gmail.com",
				Password: "password",
				Address:  "20 random St",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Test Case 2",
			Param:          "",
			ExpectedStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			b, _ := json.Marshal(&tc.Param)
			req := httptest.NewRequest(http.MethodPut, "/v1/user", strings.NewReader(string(b)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e := echo.New()
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			ctx := log.InitParams(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))

			e.Binder = NewBinder()
			e.Validator = NewValidator()

			handlerMock := &mock.HandlerMock{}
			handlerMock.
				On("Update", ctx, testifymock.Anything).
				Return(tc.Param, tc.HandlerErr)

			ctrl := NewController(handlerMock.Update, http.StatusOK, new(domain.UpdateRequest))
			err := ctrl.Handle(c)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.Equal(t, tc.ExpectedStatus, rec.Code)
			assert.NotNil(t, rec.Body)
			assert.NoError(t, err)
		})
	}
}

func TestHandleListUser(t *testing.T) {
	testCases := []struct {
		Name           string
		Param          interface{}
		HandlerResult  []domain.UserNoPass
		HandlerErr     error
		ExpectedError  error
		ExpectedStatus int
	}{
		{
			Name: "Test Case 1",
			HandlerResult: []domain.UserNoPass{{
				ID:      123,
				Name:    "John",
				Age:     18,
				Email:   "john@gmail.com",
				Address: "20 random St",
			}},
			ExpectedStatus: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)

			e := echo.New()
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			ctx := log.InitParams(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))

			e.Binder = NewBinder()
			e.Validator = NewValidator()

			handlerMock := &mock.HandlerMock{}
			handlerMock.
				On("List", ctx, testifymock.Anything).
				Return(tc.HandlerResult, tc.HandlerErr)

			ctrl := NewController(handlerMock.List, http.StatusOK, nil)
			err := ctrl.Handle(c)
			if tc.ExpectedError != nil {
				assert.Error(t, err)
				responseError := err.(*domain.ResponseError)
				assert.Equal(t, tc.ExpectedError.Error(), responseError.DeveloperMessage)
				return
			}

			assert.Equal(t, tc.ExpectedStatus, rec.Code)
			assert.NotNil(t, rec.Body)
			assert.NoError(t, err)
		})
	}
}
