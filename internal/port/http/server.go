package http

import (
	"log/slog"
	"net/http"

	_ "github.com/pdkonovalov/hitalent-test-assignment-2026/docs"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app"

	http_swagger "github.com/swaggo/http-swagger/v2"
)

func New(
	log *slog.Logger,
	app *app.App,
) *http.ServeMux {
	apiMux := http.NewServeMux()

	apiMux.Handle("POST /departments/", CreateDepartment(app.Commands.CreateDepartment))
	apiMux.Handle("POST /departments/{id}/employees/", CreateEmployee(app.Commands.CreateEmployee))
	apiMux.Handle("GET /departments/{id}", GetDepartment(app.Queries.GetDepartment))
	apiMux.Handle("PATCH /departments/{id}", UpdateDepartment(app.Commands.UpdateDepartment))
	apiMux.Handle("DELETE /departments/{id}", DeleteDepartment(app.Commands.DeleteDepartment))

	mux := http.NewServeMux()

	mux.Handle("/departments/", loggerMiddleware(log, apiMux))

	mux.Handle("GET /swagger/", http_swagger.WrapHandler)

	return mux
}
