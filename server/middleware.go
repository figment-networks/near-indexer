package server

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/gin-gonic/gin"
)

func RollbarMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				config.TrackPanic(err)
				panic(err) // continue with default panic loger
			}
		}()
		c.Next()
	}
}
