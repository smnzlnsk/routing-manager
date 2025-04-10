package handler

import (
	"encoding/json"
	"net/http"

	"github.com/smnzlnsk/routing-manager/internal/api/v1/response"
	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/service"
	"go.uber.org/zap"
)

type AlertHandler struct {
	service service.AlertService
	logger  *zap.Logger
}

func NewAlertHandler(service service.AlertService, logger *zap.Logger) *AlertHandler {
	return &AlertHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AlertHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	var req domain.Alert

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info("Handling alert", zap.Any("request", req))

	err := h.service.HandleAlert(r.Context(), &domain.Alert{
		AppName: req.AppName,
	})
	if err != nil {
		h.logger.Error("Error handling alert", zap.Error(err))
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, nil, http.StatusOK)
}
