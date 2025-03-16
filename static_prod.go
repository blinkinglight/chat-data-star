package chat

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	hashFS "github.com/benbjohnson/hashfs"
)

//go:embed static/*
var staticFS embed.FS
var staticRootFS, _ = fs.Sub(staticFS, "web/static")

func StaticProd(logger *slog.Logger) http.Handler {
	logger.Debug("static assets are embedded")
	return hashFS.FileServer(staticRootFS)
}
