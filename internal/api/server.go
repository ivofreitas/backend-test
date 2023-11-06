package api

import (
	"backend-test/internal/adapter/context"
	"backend-test/internal/adapter/log"
	"backend-test/internal/api/middleware"
	"backend-test/internal/config"
	"backend-test/internal/domain"
	"backend-test/internal/repository"
	gocontext "context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	echo   *echo.Echo
	db     *sql.DB
	logger *logrus.Entry
	signal chan struct{}
}

// NewServer - creates a new instance of server
func NewServer() *Server {
	log.Init()

	return &Server{
		logger: log.NewEntry(),
		signal: make(chan struct{}),
	}
}

// Run - start the server
func (s *Server) Run() {
	s.start()
	s.logger.Println("Server started and waiting for the graceful signal...")
	<-s.signal
}

func (s *Server) start() {
	go s.watchStop()

	serverConfig := config.GetEnv().Server

	s.logger.Infof("Server is starting in port %s.", serverConfig.Port)

	s.initEcho()

	s.db = repository.GetConn()

	Register(s.echo.Group("/v1"), Option{DB: s.db})

	addr := fmt.Sprintf(":%s", serverConfig.Port)
	go func() {
		if err := s.echo.Start(addr); err != nil {
			s.logger.WithError(err).Fatal("Shutting down the server now")
		}
	}()
}

func (s *Server) watchStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	s.stop()
}

func (s *Server) stop() {
	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Second)
	defer cancel()

	s.logger.Info("Server is stopping...")

	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.logger.Errorln(err)
	}

	err = s.db.Close()
	if err != nil {
		s.logger.Errorln(err)
	}

	close(s.signal)
}

func (s *Server) initEcho() {
	s.echo = echo.New()
	s.echo.Validator = middleware.NewValidator()
	s.echo.Binder = middleware.NewBinder()
	s.echo.Use(middleware.Logger)
	s.echo.Use(echomiddleware.Recover())
	s.echo.Pre(echomiddleware.RemoveTrailingSlash())
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		httpLog := context.Get(c.Request().Context(), log.HTTPKey).(*log.HTTP)
		httpLog.Level = logrus.WarnLevel

		responseErr := domain.ErrorDiscover(err)
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(responseErr.StatusCode)
		} else {
			err = c.JSON(responseErr.StatusCode, responseErr)
		}
		if err != nil {
			s.echo.Logger.Error(err)
		}
	}
}
