package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketplace/internal/api/advertisement"
	"marketplace/internal/api/auth"
	"marketplace/internal/api/feed"
	"marketplace/internal/api/middleware"
	"marketplace/internal/config"
	adRepo "marketplace/internal/repository/advertisement"
	userRepo "marketplace/internal/repository/user"
	adService "marketplace/internal/sevice/advertisement"
	authService "marketplace/internal/sevice/auth"
	feedService "marketplace/internal/sevice/feed"
	"marketplace/internal/validator"

	"marketplace/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	configPath = flag.String("config-path", "", "Path to config file")
	isLocal    = flag.Bool("local", false, "local run")
)

func main() {
	flag.Parse()

	mainConfig, err := config.InitMainConfig(configPath, isLocal)
	if err != nil {
		logger.Fatal("failed config loading", logger.Error, err)
	}

	aesKey := mainConfig.GetAESKey()
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := chi.NewRouter()

	logger.InitLogging()

	db, err := pgxpool.New(ctx, mainConfig.GetDbConfig())
	if err != nil {
		logger.Fatal("connect to database", logger.Error, err)
	}

	validator := validator.New()

	userRepo := userRepo.New(db)
	authService := authService.New(userRepo, mainConfig.GetJWSKey())
	authHandler := auth.New(authService, validator)

	adRepo := adRepo.New(db)
	adSevice := adService.New(adRepo)
	adHandler := advertisement.New(adSevice, validator)

	feedService := feedService.New(adRepo, aesgcm, "http://"+mainConfig.GetServerAddress())
	feedHandler := feed.New(authService, feedService, mainConfig.GetDefoultFeedLimit())

	r.Get("/sign_in", authHandler.SingIn)
	r.Get("/sign_up", authHandler.SingUp)

	r.Route("/advertisements", func(r chi.Router) {
		r.Use(
			middleware.LoggingMiddleware,
			middleware.AuthStrictMiddleware(authService),
		)
		r.Post("/create", adHandler.CreateAd)
		r.Get("/get", adHandler.GetAd)
		r.Patch("/update", adHandler.UpdateAd)
		r.Delete("/delete", adHandler.DeleteAd)
	})

	r.With(middleware.LoggingMiddleware, middleware.AuthOptionalMiddleware(authService)).
		Get("/feed", feedHandler.Feed)

	server := &http.Server{
		Addr:    mainConfig.GetServerAddress(),
		Handler: r,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("shutting down gracefully...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Fatal("HTTP server Shutdown error", logger.Error, err)
		}
	}()

	slog.Info(fmt.Sprintf("server started: %s", mainConfig.GetServerAddress()))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("HTTP server ListenAndServe error", logger.Error, err)
	}

	slog.Info("server stopped")
}
