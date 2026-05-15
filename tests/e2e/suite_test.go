package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	postgres_adapter "github.com/pdkonovalov/hitalent-test-assignment-2026/internal/adapter/postgres"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app"
	http_server "github.com/pdkonovalov/hitalent-test-assignment-2026/internal/port/http"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/migrations"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type E2ESuite struct {
	suite.Suite
	container testcontainers.Container
	handler   http.Handler
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

func (s *E2ESuite) SetupSuite() {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	s.Require().NoError(err)
	s.container = container

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Discard,
		TranslateError: true,
	})
	s.Require().NoError(err)

	s.runMigrations(db)

	slogNoop := slog.New(slog.DiscardHandler)

	txManager := postgres_adapter.NewTxManager(db)
	repository := postgres_adapter.NewRepository(db)
	readModel := postgres_adapter.NewReadModel(db)

	app := app.New(txManager, repository, repository, readModel)

	s.handler = http_server.New(slogNoop, app)
}

func (s *E2ESuite) TearDownSuite() {
	ctx := context.Background()
	if s.container != nil {
		s.container.Terminate(ctx)
	}
}

func (s *E2ESuite) runMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	s.Require().NoError(err)

	goose.SetBaseFS(migrations.FS)
	s.Require().NoError(goose.SetDialect("postgres"))
	s.Require().NoError(goose.Up(sqlDB, migrations.Dir))
}

func (s *E2ESuite) do(method, path string, body any) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		s.Require().NoError(err)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&v))
	return v
}

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

type CreateDepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name,omitempty"`
	ParentID *int64  `json:"parent_id,omitempty"`
}

type UpdateDepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at,omitempty"`
}

type CreateEmployeeResponse struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Employee struct {
	ID        int64      `json:"id"`
	FullName  string     `json:"full_name"`
	Position  string     `json:"position"`
	HiredAt   *time.Time `json:"hired_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Department struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type GetDepartmentResponse struct {
	Department *Department   `json:"department"`
	Employees  []*Employee   `json:"employees"`
	Children   []*Department `json:"children"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
