package middleware

import (
	"backend-test/internal/domain"
	"github.com/labstack/echo/v4"
)

type Binder struct{}

func NewBinder() *Binder {
	return &Binder{}
}

// Bind - information sent though request into internal structs
func (cb *Binder) Bind(i interface{}, c echo.Context) error {
	db := new(echo.DefaultBinder)

	if err := db.Bind(i, c); err != nil {
		return domain.ErrorDiscover(domain.BadRequest{DeveloperMessage: err.Error()})
	}
	return nil
}
