package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func loggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStarted := time.Now()

		requestLogCtx := newRequestLogContext(r)

		withRequestLogContext(r, requestLogCtx)

		next.ServeHTTP(w, r)

		requestLogCtx, ok := logContextFromRequest(r)
		if !ok {
			log.Error("request log context not found")

			return
		}

		requestLogCtx.DurationMs = new(time.Since(requestStarted).Milliseconds())

		if requestLogCtx.InternalServerError == true {
			log.Error("request failed", slog.Any("request", requestLogCtx))

			return
		}

		log.Info("request processed", slog.Any("request", requestLogCtx))
	})
}

type requestLogContextKey struct{}

type requestLogContext struct {
	Method              string
	Path                string
	Query               string
	Status              *string
	DurationMs          *int64
	Errors              []string
	InternalServerError bool
}

func (r requestLogContext) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0)

	attrs = append(attrs,
		slog.String("method", r.Method),
		slog.String("path", r.Path),
		slog.String("query", r.Query),
	)

	if r.Status != nil {
		attrs = append(attrs, slog.String("status", *r.Status))
	}

	if r.DurationMs != nil {
		attrs = append(attrs, slog.Int64("duration_ms", *r.DurationMs))
	}

	if len(r.Errors) != 0 {
		attrs = append(attrs, slog.Any("errors", r.Errors))
	}

	return slog.GroupValue(attrs...)
}

func newRequestLogContext(r *http.Request) requestLogContext {
	return requestLogContext{
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.Query().Encode(),
		Errors: make([]string, 0),
	}
}

func withRequestLogContext(r *http.Request, logCtx requestLogContext) {
	ctx := context.WithValue(r.Context(), requestLogContextKey{}, logCtx)

	*r = *r.WithContext(ctx)
}

func logContextFromRequest(r *http.Request) (requestLogContext, bool) {
	ctx := r.Context()

	log, ok := ctx.Value(requestLogContextKey{}).(requestLogContext)

	return log, ok
}

func withRequestLogContextStatus(r *http.Request, status int) {
	logCtx, ok := logContextFromRequest(r)
	if !ok {
		return
	}

	logCtx.Status = new(http.StatusText(status))

	if status == http.StatusInternalServerError {
		logCtx.InternalServerError = true
	}

	withRequestLogContext(r, logCtx)
}

func withRequestLogContextError(r *http.Request, err error) {
	if err == nil {
		return
	}

	logCtx, ok := logContextFromRequest(r)
	if !ok {
		return
	}

	logCtx.Errors = append(logCtx.Errors, err.Error())

	withRequestLogContext(r, logCtx)
}
