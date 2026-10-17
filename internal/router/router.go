package router

import (
	"time"

	handler "stream/internal/handler/stream"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(streamHandler *handler.StreamHandler) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://192.168.29.87:3000",
			"https://192.168.29.87:3000",
			"https://video.dragskills.in",
			"https://cac1-2405-201-2032-d09b-415e-d58b-a0d-ff2a.ngrok-free.app",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"ngrok-skip-browser-warning",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	streams := api.Group("/streams")

	streams.POST("", streamHandler.CreateStream)
	streams.GET("", streamHandler.ListStreams)
	streams.GET("/:id", streamHandler.GetStreamByID)
	streams.POST("/:id/start", streamHandler.StartStream)
	streams.POST("/:id/token", streamHandler.GenerateToken)
	streams.POST("/:id/end", streamHandler.EndStream)

	return r
}
