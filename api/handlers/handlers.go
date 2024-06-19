package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/macadrich/go-bike/api"
	_ "github.com/macadrich/go-bike/docs"
	"github.com/macadrich/go-bike/pkg/utils"
)

type Handlers struct {
	svc api.IService
}

func NewHandlers(svc api.IService) *Handlers {
	return &Handlers{svc}
}

// HealthCheck godoc
// @Summary Show if api is running
// @Description get string health check
// @Tags example
// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthcheck [get]
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "health check ok!")
}

func (h *Handlers) InsertStation(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.InsertStation(r.Context()); err != nil {
		sendResponse(w, http.StatusServiceUnavailable, ErrorMessage{
			Message: "Unable to update stations",
		})
		return
	}

	sendResponse(w, http.StatusOK, ResponseMessage{
		Message: "Update station successfully!",
	})
}

func (h *Handlers) QueryAllStation(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("at")

	if !utils.ValidateTimestamp(status) {
		sendResponse(w, http.StatusBadRequest, ErrorMessage{
			Message: "condition is invalid",
		})
		return
	}

	result, err := h.svc.QueryAllStation(r.Context(), status)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, ErrorMessage{
			Message: "unable to get stations",
		})
		return
	}

	sendResponse(w, http.StatusOK, result)
}

func (h *Handlers) QuerySpecificStation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "kioskId")
	status := r.URL.Query().Get("at")

	if !utils.ValidateTimestamp(status) {
		sendResponse(w, http.StatusBadRequest, ErrorMessage{
			Message: "condition is invalid",
		})
		return
	}

	kioskId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	station, err := h.svc.QuerySpecificStation(kioskId, status)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, ErrorMessage{
			Message: "Unable to get station",
		})
		return
	}

	sendResponse(w, http.StatusOK, station)
}
