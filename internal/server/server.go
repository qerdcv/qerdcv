package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/qerdcv/qerdcv/internal/config"
	"github.com/qerdcv/qerdcv/internal/server/handlers"
	"github.com/qerdcv/qerdcv/internal/server/middlewares"
	"github.com/qerdcv/qerdcv/web"
)

type Server struct {
	server *http.Server
	logger *slog.Logger

	shutdownTimeout time.Duration
}

func New(
	logger *slog.Logger,
	cfg config.ServerConfig,
	authMiddleware echo.MiddlewareFunc,
	userHandler *handlers.UserHandler,
	budgetHandler *handlers.BudgetHandler,
) *Server {
	e := echo.New()
	e.Renderer = NewTemplateRenderer()
	e.Use(
		middlewares.Recover(logger),
		middlewares.Logging(logger),
	)

	// ====static handling====
	staticHandler := echo.StaticDirectoryHandler(web.Static, false)
	e.GET("/static/*",
		func(c echo.Context) error {
			c.SetParamValues("/static/" + c.Param("*"))
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int((24*time.Hour).Seconds())))
			return staticHandler(c)
		},
	)
	// ======================

	e.GET("", func(c echo.Context) error {
		return c.Render(http.StatusOK, "templates/index.gohtml", nil)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(http.StatusOK, "templates/auth/index.gohtml", nil)
	})

	apiG := e.Group("/api/v1")

	registerUserRoutes(apiG.Group("/users"), userHandler)
	registerBudgetRoutes(apiG.Group("/budget", authMiddleware), budgetHandler)

	return &Server{
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      e,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		logger:          logger,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		defer s.server.Close()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed to shutdown server: %w", err)
			return
		}
	}()

	s.logger.Info("starting up server")
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server listen and serve: %w", err)
	}

	return nil
}

func registerUserRoutes(g *echo.Group, handler *handlers.UserHandler) {
	g.POST("", handler.CreateUser)
	g.POST("/auth", handler.AuthorizeUser)
}

func registerBudgetRoutes(g *echo.Group, handler *handlers.BudgetHandler) {
	g.POST("/categories", handler.CreateCategory)
	g.GET("/categories", handler.CategoriesList)

	g.POST("/transactions", handler.CreateTransaction)
	g.GET("/transactions", handler.TransactionsList)
}
