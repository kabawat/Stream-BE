package handler

import (
	"context"
	"errors"
	"net/http"
	"stream/internal/repository"
	"stream/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GenerateTokenRequest struct {
	Role string `json:"role" binding:"required,oneof=host viewer"`
}

func (h *StreamHandler) GenerateToken(c *gin.Context) {
	var req GenerateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "role must be either host or viewer",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	stream, err := h.service.GetStreamByID(
		ctx,
		c.Param("id"),
	)
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

	// Host can obtain a token only after start has been requested.
	if req.Role == "host" && stream.Status != "starting" {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Stream must be starting before host can connect",
		})
		return
	}

	if req.Role == "viewer" &&
		stream.Status != "starting" &&
		stream.Status != "live" {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Stream is not available to join",
		})
		return
	}

	// Generate a server-side participant identity.
	identity := bson.NewObjectID().Hex()

	token, err := h.livekitService.GenerateToken(
		stream.RoomName,
		req.Role,
		identity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate LiveKit token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"serverUrl": streamServerURL(h.livekitService),
		"token":     token,
		"identity":  identity,
		"roomName":  stream.RoomName,
		"role":      req.Role,
	})
}

func streamServerURL(livekitService *service.LiveKitService) string {
	return livekitService.ServerURL()
}
