package app

import (
	"context"
	"log"
	"net/http"
	"os"
)

// Run starts app by running http server
func (app *App) Run(done chan os.Signal) {
	app.IDMappers.RunReloader(app.log)
	defer app.IDMappers.StopReloader()

	httpServer := &http.Server{
		Addr:    app.Configuration.Addr,
		Handler: CreateRouter(app.Configuration.APIPrefix, app.Version, app.IDMappers),
	}

	go func() {
		app.log.Infof("Starting server on %s", app.Configuration.Addr)
		app.log.Fatal(httpServer.ListenAndServe())
	}()

	<-done
	app.log.Info("Shutdown signal received. Exiting.")

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to shutdown http server: %s", err)
	}
}
