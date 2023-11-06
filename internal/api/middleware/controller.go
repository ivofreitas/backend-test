package middleware

import (
	"backend-test/internal/adapter/context"
	"backend-test/internal/adapter/log"
	"backend-test/internal/domain"
	gocontext "context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type HandlerFn func(ctx gocontext.Context, param interface{}) (interface{}, error)

type controller struct {
	fn         HandlerFn
	param      interface{}
	httpStatus int
}

func NewController(fn HandlerFn, httpStatus int, param interface{}) *controller {
	return &controller{fn, param, httpStatus}
}

// Handle - Validate the request
func (ctrl *controller) Handle(c echo.Context) error {

	ctx := c.Request().Context()
	httpLog := context.Get(ctx, log.HTTPKey).(*log.HTTP)

	if ctrl.param != nil {
		ctrl.param = reflect.New(reflect.TypeOf(ctrl.param).Elem()).Interface()
		if err := c.Bind(ctrl.param); err != nil {
			responseErr := err.(*domain.ResponseError)
			httpLog.Error = responseErr.Error()
			return c.JSON(http.StatusBadRequest, responseErr)
		}

		if err := c.Validate(ctrl.param); err != nil {
			responseErr := err.(*domain.ResponseError)
			httpLog.Error = responseErr.Error()
			return c.JSON(http.StatusBadRequest, responseErr)
		}

		b, _ := json.Marshal(ctrl.param)
		httpLog.Request.Param = string(b)
	}

	result, err := ctrl.fn(ctx, ctrl.param)
	if err != nil {
		responseErr := err.(*domain.ResponseError)
		httpLog.Error = responseErr.Error()
		return c.JSON(responseErr.StatusCode, responseErr)
	}

	httpLog.Level = logrus.InfoLevel
	if result != nil {
		httpLog.Response.Body = result
		return c.JSON(ctrl.httpStatus, result)
	}

	return c.JSON(http.StatusNoContent, nil)
}
