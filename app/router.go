package app

import (
	"fmt"

	"github.com/danielkraic/idmapper/app/handlers"
	"github.com/danielkraic/idmapper/app/idmappers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CreateRouter creates http router
func CreateRouter(apiPrefix string, appVersion *handlers.Version, idMappers *idmappers.IDMappers) *mux.Router {
	r := mux.NewRouter()

	versioned := func(route string) string {
		return fmt.Sprintf("%s%s", apiPrefix, route)
	}

	r.Handle(versioned("/currency/{id}"), handlers.NewIDMapperHandler(idMappers.CurrencyCodes)).Methods("GET")
	r.Handle(versioned("/country/{id}"), handlers.NewIDMapperHandler(idMappers.CountryCodes)).Methods("GET")
	r.Handle(versioned("/language/{id}"), handlers.NewIDMapperHandler(idMappers.LanguageCodes)).Methods("GET")

	r.Handle("/version", handlers.NewVersionHandler(appVersion)).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandlerFunc).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	return r
}
