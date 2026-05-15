package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgres_adapter "github.com/pdkonovalov/hitalent-test-assignment-2026/internal/adapter/postgres"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/config"
	http_server "github.com/pdkonovalov/hitalent-test-assignment-2026/internal/port/http"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/migrations"
	"github.com/pressly/goose/v3"

	swagger_docs "github.com/pdkonovalov/hitalent-test-assignment-2026/docs"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

// @title		Departments API
// @version	1.0
// @host		localhost:8080
// @BasePath	/
func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := newLogger(cfg)

	db, err := newGorm(cfg)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}

	if err := doMigrationsUp(db); err != nil {
		return fmt.Errorf("do migrations up: %w", err)
	}

	setupSwaggerDocs(cfg)

	txManager := postgres_adapter.NewTxManager(db)
	repository := postgres_adapter.NewRepository(db)
	readModel := postgres_adapter.NewReadModel(db)

	app := app.New(txManager, repository, repository, readModel)

	api := http_server.New(log, app)

	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: http.StripPrefix(cfg.BasePath, api),
	}

	go func() {
		log.Info("starting server",
			slog.String("host", cfg.Host),
			slog.String("port", cfg.Port),
			slog.String("base_path", cfg.BasePath),
		)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server stopped")

	return nil
}

func newLogger(cfg *config.Config) *slog.Logger {
	opts := slog.HandlerOptions{
		Level: cfg.LogLevel.SlogLevel(),
	}

	if cfg.LogFormat == config.LogFormatText {
		return slog.New(slog.NewTextHandler(os.Stdout, &opts))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &opts))
}

func newGorm(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	return gorm.Open(
		gorm_postgres.Open(dsn),
		&gorm.Config{
			Logger:         gorm_logger.Default.LogMode(gorm_logger.Silent),
			TranslateError: true,
		},
	)
}

func doMigrationsUp(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(sqlDB, migrations.Dir)
}

func setupSwaggerDocs(cfg *config.Config) {
	swagger_docs.SwaggerInfo.Host = cfg.Endpoint
	swagger_docs.SwaggerInfo.BasePath = cfg.BasePath
}
