package api

import (
	"backend-test/internal/api/handler"
	"backend-test/internal/api/middleware"
	"backend-test/internal/config"
	"backend-test/internal/domain"
	"backend-test/internal/repository"
	"backend-test/internal/swagger"
	"database/sql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Option struct {
	DB *sql.DB
}

// Register - create all routes
func Register(g *echo.Group, opts Option) {

	env := config.GetEnv()
	doc := env.Doc
	if doc.Enabled {
		swagger.Register(swagger.Options{
			Title:       doc.Title,
			Description: doc.Description,
			Version:     doc.Version,
			BasePath:    env.Server.BasePath,
			Group:       g.Group("/swagger"),
		})
	}

	g.GET("/health", handler.Health)

	userRoute(g, opts)
}

func userRoute(g *echo.Group, opts Option) {
	env := config.GetEnv().MySQL

	userRepo := repository.NewUserRepo(env.Timeout, opts.DB)
	userHdl := handler.NewUserHdl(userRepo, bcrypt.GenerateFromPassword)
	crtUsr := middleware.NewController(userHdl.Create, http.StatusCreated, new(domain.User))
	getUsr := middleware.NewController(userHdl.GetByID, http.StatusOK, new(domain.GetByIDRequest))
	listUsr := middleware.NewController(userHdl.List, http.StatusOK, nil)
	uptUsr := middleware.NewController(userHdl.Update, http.StatusOK, new(domain.UpdateRequest))
	dltUsr := middleware.NewController(userHdl.Delete, http.StatusNoContent, new(domain.DeleteRequest))

	userGroup := g.Group("/user")
	userGroup.POST("", crtUsr.Handle)
	userGroup.GET("/:id", getUsr.Handle)
	userGroup.GET("", listUsr.Handle)
	userGroup.PUT("/:id", uptUsr.Handle)
	userGroup.DELETE("/:id", dltUsr.Handle)
}
