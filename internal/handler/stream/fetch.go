package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StreamHandler) ListStreams(c *gin.Context) {
	limit := int64(20)

	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil || parsedLimit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be a positive integer",
			})
			return
		}

		limit = parsedLimit
	}

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	streams, err := h.service.ListStreams(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch streams",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(streams),
		"data":  streams,
	})
}
