package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/config"
)

// SetGinDefaults changes Gin behavior base on application environment
func SetGinDefaults(cfg *config.Config) {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}
