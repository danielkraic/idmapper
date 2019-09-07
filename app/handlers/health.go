package handlers

import "net/http"

// HealthHandlerFunc handles /health and return http 200 status code
func HealthHandlerFunc(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
