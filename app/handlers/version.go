package handlers

import (
	"encoding/json"
	"net/http"
)

// Version is application version details
type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Build   string `json:"build"`
}

type versionHandler struct {
	version *Version
}

// NewVersionHandler creates Version handler
func NewVersionHandler(version *Version) http.Handler {
	return &versionHandler{
		version: version,
	}
}

// ServeHTTP serves http request
func (h *versionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(h.version)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
