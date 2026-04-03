package handlers

import (
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	configService services.ConfigService
}

func NewConfigHandler(service services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: service,
	}
}

func (h *ConfigHandler) GetConfig(c echo.Context) error {
	config := h.configService.GetAppConfig()
	return c.JSON(http.StatusOK, config)
}
