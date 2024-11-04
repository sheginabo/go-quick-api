package handlers

import (
	"encoding/json"
	"net/http"
)

type InternalHandler struct{}

func NewInternalHandler() *InternalHandler {
	return &InternalHandler{}
}

func (h *InternalHandler) PostHello(w http.ResponseWriter, r *http.Request) {
	var req PostHelloRequest
	if err := ValidatePayload(r, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ToErrorResponse(*err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello " + req.Message,
	})
}
