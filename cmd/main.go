package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/instaUpload/user-service/server"
	t "github.com/instaUpload/user-service/types"
	u "github.com/instaUpload/user-service/utils"
)

type App struct {
	Server server.Server
	Config *t.ApplicationConfig
}

func (a *App) Run(_ context.Context) error {
	// Start the server
	a.Server.MountRoutes()
	return http.ListenAndServe(a.Config.Address(), a.Server.Handler())
}

func main() {
	// STEP 1: Read all the environment variables needed for the application.
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	// STEP 2: Initialize the logger.
	app_env := u.GetEnvAsString("APP_ENV", "development")
	options := &slog.HandlerOptions{
		AddSource: app_env != "production",
		Level: func() slog.Level {
			if app_env == "production" {
				return slog.LevelInfo
			}
			return slog.LevelDebug
		}(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(a.Value.String()[len(cwd):])
			}
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.String()[:19] + " UTC")
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, options))
	slog.SetDefault(logger)
	slog.Info("Logger initialized", "app_env", app_env)

	// STEP 3: Create application level context with cancel function.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Step 4: Setup the application.
	appConfig := t.NewApplicationConfig()
	app := App{
		Config: appConfig,
		Server: server.NewChiServer(ctx),
	}

	// Step 5: Start the application in a goroutine.
	slog.Info("Application started...")
	fmt.Printf("Server running at https://%s\n", appConfig.Address())
	go app.Run(ctx)

	// Step 6: Listen for OS signals to gracefully shutdown the application.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	// Waiting for a signal
	<-sig
	slog.Warn("Shutting down application...")
	cancel()
	time.Sleep(2 * time.Second) // Simulate some cleanup work
	<-ctx.Done()
	slog.Warn("Application stopped.")

}
