package main

import (
	"context"
	"fmt"
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	app, err := App()
	if err != nil {
		return err
	}
	cfg := app.Config

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Http.Port),
		Handler: app.GinEngine.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.ZapLogger.Warnw("server is shutting down in 5 seconds", "event", "SERVER_SHUTDOWN")
	if err := server.Shutdown(ctx); err != nil {
		app.ZapLogger.Error(err)
	}

	return nil
}
