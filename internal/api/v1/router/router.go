package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/smnzlnsk/routing-manager/internal/api/v1/handler"
	"github.com/smnzlnsk/routing-manager/internal/service"
	"go.uber.org/zap"
)

func Setup(services *service.Services, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)

	interestHandler := handler.NewInterestHandler(services.InterestService, logger)

	router.Route("/api/v1/interests", func(r chi.Router) {
		r.Post("/", interestHandler.Create)
		r.Get("/", interestHandler.List)

		r.Get("/app/{appName}", interestHandler.GetByAppName)
		r.Delete("/app/{appName}", interestHandler.DeleteByAppName)

		r.Get("/service/{serviceIp}", interestHandler.GetByServiceIp)
		r.Delete("/service/{serviceIp}", interestHandler.DeleteByServiceIp)
	})

	alertHandler := handler.NewAlertHandler(services.AlertService, logger)

	router.Route("/api/v1/alert", func(r chi.Router) {
		r.Post("/", alertHandler.HandleAlert)

	})

	routingHandler := handler.NewRoutingHandler(services.RoutingService, logger)
	router.Route("/api/v1/routing", func(r chi.Router) {
		r.Post("/", routingHandler.HandleRoutingChange)
		r.Get("/app/{appName}", routingHandler.GetRouting)
	})

	return router
}
