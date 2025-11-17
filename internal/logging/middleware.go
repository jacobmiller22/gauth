package logging

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobmiller22/gauth/pkg/clog"
)

func LogReq(l *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l = l.With("request-id", uuid.NewString())
		r.WithContext(clog.WithContext(r.Context(), l))

		l.InfoContext(r.Context(), "request-received", "method", r.Method, "path", r.URL.String())

		h.ServeHTTP(w, r) // Call the next handler in the chain
	})
}
