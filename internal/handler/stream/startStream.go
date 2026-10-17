package handler

import (
	"context"
	"errors"
	"net/http"
	"stream/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StreamHandler) StartStream(c *gin.Context) {
	streamID := c.Param("id")

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	stream, err := h.service.StartStream(ctx, streamID)
	if err != nil {
		switch {
		case err.Error() == "invalid stream ID":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, repository.ErrStreamNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Stream not found",
			})

		case errors.Is(err, repository.ErrStreamConflict):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Stream cannot be started in its current status",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to start stream",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stream start initiated successfully",
		"data":    stream,
	})
}
