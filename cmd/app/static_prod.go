//go:build !dev
// +build !dev

package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/blinkinglight/chat-data-star"
)

func static() http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Static assets are being served from /static/")
	return chat.StaticProd(logger)
}
