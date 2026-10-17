package handler

import (
	"context"
	"net/http"
	"stream/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StreamHandler) CreateStream(c *gin.Context) {
	var input service.CreateStreamInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	stream, err := h.service.CreateStream(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stream created successfully",
		"data":    stream,
	})
}
