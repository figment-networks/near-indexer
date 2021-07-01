package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

// jsonError renders an error response
func jsonError(c *gin.Context, status int, err interface{}) {
	var message interface{}

	switch v := err.(type) {
	case error:
		message = v.Error()
	default:
		message = v
	}

	c.AbortWithStatusJSON(status, gin.H{
		"status": status,
		"error":  message,
	})
}

// badRequest renders a HTTP 400 bad request response
func badRequest(c *gin.Context, err interface{}) {
	jsonError(c, http.StatusBadRequest, err)
}

// notFound renders a HTTP 404 not found response
func notFound(c *gin.Context, err interface{}) {
	jsonError(c, http.StatusNotFound, err)
}

// serverError renders a HTTP 500 error response
func serverError(c *gin.Context, err interface{}) {
	jsonError(c, http.StatusInternalServerError, err)
}

// jsonOk renders a successful response
func jsonOk(c *gin.Context, data interface{}) {
	switch data.(type) {
	case []byte:
		c.Header("Content-Type", "application/json")
		c.String(200, "%s", data)
	default:
		c.JSON(200, data)
	}
}

// shouldReturn is a shorthand method for handling resource errors
func shouldReturn(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if err == store.ErrNotFound || err == near.ErrNotExist {
		notFound(c, err)
	} else {
		serverError(c, err)
	}

	return true
}
