package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kube-tarian/container-bridge/agent/pkg/clients"
	"github.com/kube-tarian/container-bridge/agent/pkg/config"
	"github.com/kube-tarian/container-bridge/agent/pkg/publish"

	"github.com/julienschmidt/httprouter"
)

type Application struct {
	Config  *config.Config
	server  *http.Server
	conn    *clients.NATSContext
	Publish publish.Models
}

func New(conf *config.Config, conn *clients.NATSContext) *Application {
	app := &Application{
		Config: conf,
		conn:   conn,
	}

	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return app
}

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/localregistry/event", app.localRegistryHandler)
	return router
}

func (app *Application) Start() {
	log.Println("Starting server on port", app.Config.Port)
	if err := app.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server closed, readon: %v", err)
	}
}

func (app *Application) Close() {
	log.Printf("Closing the service gracefully")
	app.conn.Close()

	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Printf("Could not close the service gracefully: %v", err)
	}
}
