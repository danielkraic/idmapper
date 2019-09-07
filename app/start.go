package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielkraic/idmapper/app/idmappers"
)

// Run starts app by running http server
func (app *App) Run() {
	idMappers, err := idmappers.NewIDMappers(app.log, app.redisClient, app.db, &app.configuration.IDMappers)
	if err != nil {
		app.log.Fatalf("failed to create IDMappers: %s", err)
	}
	idMappers.RunReloader(app.log)

	httpServer := &http.Server{
		Addr:    app.configuration.Addr,
		Handler: createRouter(app.configuration.APIPrefix, app.version, idMappers),
	}

	go func() {
		app.log.Infof("Starting server on %s", app.configuration.Addr)
		app.log.Fatal(httpServer.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	app.log.Info("Shutdown signal received. Exiting.")

	err = httpServer.Shutdown(context.Background())
	idMappers.StopReloader()

	if err != nil {
		log.Fatalf("failed to shutdown http server: %s", err)
	}
}
