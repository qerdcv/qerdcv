package http

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"
)

const staticMaxAge = 24 * time.Hour

func NewHandler(
	logger *slog.Logger,
	staticFS fs.FS,
	profileH *ProfileHandler,
) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", profileH.Index)

	fileServer := http.FileServer(http.FS(staticFS))
	m.Handle(
		"/static/",
		http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// wrapper for cache control headers
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int(staticMaxAge.Seconds())))
			fileServer.ServeHTTP(w, r)
		})))
	return LoggingMiddleware(logger)(m)
}
