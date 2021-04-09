package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fberrez/simple-uptime-backend/api"
	"github.com/fberrez/simple-uptime-backend/backend/postgres"
	"github.com/gin-gonic/gin"
	"github.com/ovh/configstore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// portKey is the key corresponding to the port value in the config file
var portKey = "port"

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
}

func main() {
	// initializes config variables
	configstore.InitFromEnvironment()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// initializes backend
	back, err := postgres.New()
	if err != nil {
		panic(err)
	}

	defer back.Close()

	a := api.New(back)

	port, err := configstore.GetItemValueInt(portKey)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Cannot start server")
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
