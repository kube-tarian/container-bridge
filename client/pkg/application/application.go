package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kube-tarian/container-bridge/client/pkg/clickhouse"
	"github.com/kube-tarian/container-bridge/client/pkg/clients"
	"github.com/kube-tarian/container-bridge/client/pkg/config"
	"github.com/kube-tarian/container-bridge/client/pkg/handler"
)

type Application struct {
	Config     *config.Config
	apiServer  *handler.APIHandler
	httpServer *http.Server
	conn       *clients.NATSContext
	dbClient   *clickhouse.DBClient
}

func New(conf *config.Config, conn *clients.NATSContext, dbClient *clickhouse.DBClient) *Application {
	log.Println("Initializing Application")
	apiServer, err := handler.NewAPIHandler(conn)
	if err != nil {
		log.Fatalf("API Handler initialisation failed: %v", err)
	}

	mux := chi.NewMux()
	apiServer.BindRequest(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.Port),
		Handler: mux,
	}

	return &Application{
		Config:     conf,
		conn:       conn,
		dbClient:   dbClient,
		apiServer:  apiServer,
		httpServer: httpServer,
	}
}

func (app *Application) Start() {
	log.Println("Starting server on port", app.httpServer.Addr)
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Unexpected server close: %v", err)
	}
}

func (app *Application) Close() {
	log.Printf("Closing the service gracefully")
	app.conn.Close()

	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Could not close the service gracefully: %v", err)
	}
}
