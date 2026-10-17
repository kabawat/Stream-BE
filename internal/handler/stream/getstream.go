package handler

import (
	"context"
	"errors"
	"net/http"
	"stream/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StreamHandler) GetStreamByID(c *gin.Context) {
	streamID := c.Param("id")

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	stream, err := h.service.GetStreamByID(ctx, streamID)

	if err != nil {
		switch {
		case err.Error() == "invalid stream ID":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid stream ID",
			})

		case errors.Is(err, repository.ErrStreamNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Stream not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch stream",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stream fetched successfully",
		"data":    stream,
	})
}
