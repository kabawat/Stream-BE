package handler

import (
	"stream/internal/service"
)

type StreamHandler struct {
	service        *service.StreamService
	livekitService *service.LiveKitService
}

func NewStreamHandler(
	streamService *service.StreamService,
	livekitService *service.LiveKitService,
) *StreamHandler {
	return &StreamHandler{
		service:        streamService,
		livekitService: livekitService,
	}
}
