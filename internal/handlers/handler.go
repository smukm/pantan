package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"pantan/config"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) InitRoutes() *mux.Router {

	router := mux.NewRouter()

	// Middleware
	router.Use(h.loggingMiddleware)
	router.Use(h.recoveryMiddleware)
	router.Use(prometheusMiddleware)

	// Health check
	router.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/ready", h.ReadinessCheck).Methods(http.MethodGet)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", h.Status).Methods(http.MethodGet)
	api.HandleFunc("/hello", h.Hello).Methods(http.MethodGet)

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	return router
}

// --- Handlers ---

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "pantan",
	})
}

func (h *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"ready": "true",
	})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"service": "pantan",
		"version": config.Version,
		"env":     h.cfg.Env,
	})
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   fmt.Sprintf("Hello, %s!", name),
		"service":   "pantan",
		"timestamp": time.Now().Unix(),
	})
}

// --- Helpers ---

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response")
	}
}
