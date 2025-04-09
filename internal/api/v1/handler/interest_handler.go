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

type InterestHandler struct {
	service service.InterestService
	logger  *zap.Logger
}

func NewInterestHandler(service service.InterestService, logger *zap.Logger) *InterestHandler {
	return &InterestHandler{
		service: service,
		logger:  logger,
	}
}

func (h *InterestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.InterestRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	h.logger.Info("Creating interest", zap.Any("request", req))

	interest, err := h.service.Create(r.Context(), &domain.Interest{
		AppName:   req.AppName,
		ServiceIp: req.ServiceIp,
	})
	if err != nil {
		h.logger.Error("Error creating interest", zap.Error(err))
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, interest, http.StatusCreated)
}

func (h *InterestHandler) GetByAppName(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "appName")

	interest, err := h.service.GetByAppName(r.Context(), appName)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, interest, http.StatusOK)
}

func (h *InterestHandler) GetByServiceIp(w http.ResponseWriter, r *http.Request) {
	serviceIp := chi.URLParam(r, "serviceIp")

	interest, err := h.service.GetByServiceIp(r.Context(), serviceIp)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, interest, http.StatusOK)
}

func (h *InterestHandler) List(w http.ResponseWriter, r *http.Request) {
	interests, err := h.service.List(r.Context())
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, interests, http.StatusOK)
}

func (h *InterestHandler) DeleteByAppName(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "appName")

	err := h.service.DeleteByAppName(r.Context(), appName)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, nil, http.StatusOK)
}

func (h *InterestHandler) DeleteByServiceIp(w http.ResponseWriter, r *http.Request) {
	serviceIp := chi.URLParam(r, "serviceIp")

	err := h.service.DeleteByServiceIp(r.Context(), serviceIp)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, nil, http.StatusOK)
}
