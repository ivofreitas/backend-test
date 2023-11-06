package swagger

import (
	"backend-test/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Options
type Options struct {
	Group       *echo.Group
	AccessKey   string
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

// Register group item check
func Register(opts Options) {

	docs.SwaggerInfo.Title = opts.Title
	docs.SwaggerInfo.Description = opts.Description
	docs.SwaggerInfo.Version = opts.Version
	docs.SwaggerInfo.Host = opts.Host
	docs.SwaggerInfo.BasePath = opts.BasePath

	docs.SwaggerInfo.Schemes = []string{"http"}

	opts.Group.GET("/*", echoSwagger.WrapHandler)

}
