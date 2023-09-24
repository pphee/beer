package monitorHandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/peedans/beerleo/config"
	"github.com/peedans/beerleo/modules/monitorHandlers"
	"net/http"
)

type IMonitorHandler interface {
	HealthCheck(c *gin.Context)
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *gin.Context) {
	res := &monitorHandlers.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	c.JSON(http.StatusOK, res)
}
