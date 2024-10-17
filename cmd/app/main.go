package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/blinkinglight/chat-data-star/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

var port = ":8080"

func init() {
	if p := os.Getenv("BIND"); p != "" {
		port = p
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(fmt.Sprintf("Starting Server @%s", port))
	defer logger.Info("Stopping Server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, logger); err != nil {
		logger.Error("Error running server", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(startServer(ctx, logger, port))

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	return nil
}

func startServer(ctx context.Context, logger *slog.Logger, port string) func() error {
	return func() error {
		router := chi.NewMux()

		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		router.Handle("/*", static())

		handlers.SetupRoutes(logger, router)

		srv := &http.Server{
			Addr:    fmt.Sprintf("%s", port),
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()
	}
}
