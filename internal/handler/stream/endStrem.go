package handler

import (
	"errors"
	"net/http"
	"stream/internal/repository"

	"github.com/gin-gonic/gin"
)

func (h *StreamHandler) EndStream(c *gin.Context) {
	stream, err := h.service.EndStream(
		c.Request.Context(),
		c.Param("id"),
	)
	if err != nil {
		switch {
		case err.Error() == "invalid stream ID":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, repository.ErrStreamNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, repository.ErrStreamConflict):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to end stream",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stream ended successfully",
		"stream":  stream,
	})
}
