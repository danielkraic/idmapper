package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/gorilla/mux"
)

// IDMapperResponse response struct consists of ID and Name
type IDMapperResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type idMapperHandler struct {
	idMapper *idmapper.IDMapper
}

// NewIDMapperHandler creates new http.Handler with IDMapper
func NewIDMapperHandler(idMapper *idmapper.IDMapper) http.Handler {
	return &idMapperHandler{
		idMapper: idMapper,
	}
}

// ServeHTTP servers http requuests
func (h *idMapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	name, found := h.idMapper.Get(id)

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := IDMapperResponse{
		ID:   id,
		Name: name,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
