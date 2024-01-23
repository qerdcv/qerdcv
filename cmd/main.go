package main

import (
	"context"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/qerdcv/qerdcv/internal/config"
	"github.com/qerdcv/qerdcv/internal/http"
	"github.com/qerdcv/qerdcv/web"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{slog.String("version", config.Version)}))

	tplFS, err := template.ParseFS(web.Template, "template/*")
	if err != nil {
		logger.Error("template parse fs", slog.Any("error", err))
		return
	}

	tpl := http.NewTemplate(tplFS)
	staticFS, err := fs.Sub(web.Static, "assets")
	if err != nil {
		logger.Error("fs sub", slog.Any("error", err))

		return
	}

	s := http.NewHandler(logger, staticFS, http.NewProfileHandler(tpl))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := http.RunServer(ctx, s, logger); err != nil {
		logger.Error("http run server", slog.Any("error", err.Error()))
	}
}
