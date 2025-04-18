package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smnzlnsk/routing-manager/internal/api/v1/response"
	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/service"
	"go.uber.org/zap"
)

type RoutingHandler struct {
	service service.RoutingService
	logger  *zap.Logger
}

func NewRoutingHandler(service service.RoutingService, logger *zap.Logger) *RoutingHandler {
	return &RoutingHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RoutingHandler) HandleRoutingChange(w http.ResponseWriter, r *http.Request) {
	var req domain.RoutingChange

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info("Handling routing change", zap.Any("request", req))

	err := h.service.HandleRoutingChange(r.Context(), &domain.RoutingChange{
		AppName: req.AppName,
	})
	if err != nil {
		h.logger.Error("Error handling routing change", zap.Error(err))
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, nil, http.StatusOK)
}

func (h *RoutingHandler) GetRouting(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "appName")

	routing, err := h.service.GetRouting(r.Context(), appName)
	if err != nil {
		h.logger.Error("Error getting routing", zap.Error(err))
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, routing, http.StatusOK)
}
