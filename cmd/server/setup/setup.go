package setup

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-envconfig"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddr  string `env:"SERVER_ADDR,default=:8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	LogLevel    string `env:"LOG_LEVEL,default=info"`
}

// App holds all dependencies for the application
type App struct {
	Config *Config
	Logger *slog.Logger
	DB     *sqlx.DB
}

// NewApp creates a new application instance with all dependencies
func NewApp(ctx context.Context) (*App, error) {
	// Load configuration
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	// Setup logger
	var logLevel slog.Level
	switch config.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Setup database connection
	db, err := sqlx.Connect("pgx", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to database")

	return &App{
		Config: &config,
		Logger: logger,
		DB:     db,
	}, nil
}

// Close cleans up application resources
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
