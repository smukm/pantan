package handlers

import (
	"pantan/config"

	gin "github.com/gin-gonic/gin"
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

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
