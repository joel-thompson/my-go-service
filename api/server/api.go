package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/joel-thompson/my-go-service/storage"
)

// API holds the server dependencies
type API struct {
	logger *slog.Logger
	store  *storage.Store
}

// New creates a new API instance
func New(logger *slog.Logger, db *sqlx.DB) *API {
	return &API{
		logger: logger,
		store:  storage.New(db),
	}
}

// SetupRoutes configures all API routes
func (a *API) SetupRoutes() *gin.Engine {
	// Set Gin to release mode to reduce log verbosity
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(a.loggingMiddleware())

	// Health check endpoint
	router.GET("/health", a.handleHealth)

	// Hello world endpoint
	router.GET("/hello", a.handleHello)

	// Items endpoints
	router.POST("/items", a.handleCreateItem)
	router.GET("/items", a.handleListItems)

	return router
}

// loggingMiddleware adds structured logging to all requests
func (a *API) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		a.logger.Info("HTTP request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
		)
		c.Next()
	}
}
