package service

import (
	"context"
	"errors"
	"strings"

	"stream/internal/model"
	"stream/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StreamService struct {
	repo *repository.StreamRepository
}

func NewStreamService(
	repo *repository.StreamRepository,
) *StreamService {
	return &StreamService{
		repo: repo,
	}
}

type CreateStreamInput struct {
	Title      string `json:"title" binding:"required,min=3,max=100"`
	SourceType string `json:"sourceType" binding:"required,oneof=camera screen camera_screen"`
}

// create stream
func (s *StreamService) CreateStream(
	ctx context.Context,
	input CreateStreamInput,
) (*model.Stream, error) {
	title := strings.TrimSpace(input.Title)

	if len(title) < 3 || len(title) > 100 {
		return nil, errors.New(
			"title must be between 3 and 100 characters",
		)
	}

	streamID := bson.NewObjectID()

	stream := &model.Stream{
		ID:          streamID,
		RoomName:    "stream-" + streamID.Hex(),
		Title:       title,
		Status:      "created",
		SourceType:  input.SourceType,
		ViewerCount: 0,
	}

	if err := s.repo.Create(ctx, stream); err != nil {
		return nil, err
	}

	return stream, nil
}

// get stream
func (s *StreamService) ListStreams(
	ctx context.Context,
	limit int64,
) ([]model.Stream, error) {
	if limit < 1 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return s.repo.List(ctx, limit)
}

// get signle stream
func (s *StreamService) GetStreamByID(
	ctx context.Context,
	streamID string,
) (*model.Stream, error) {
	id, err := bson.ObjectIDFromHex(streamID)
	if err != nil {
		return nil, errors.New("invalid stream ID")
	}

	return s.repo.GetByID(ctx, id)
}

// start stream
func (s *StreamService) StartStream(
	ctx context.Context,
	streamID string,
) (*model.Stream, error) {
	id, err := bson.ObjectIDFromHex(streamID)
	if err != nil {
		return nil, errors.New("invalid stream ID")
	}

	return s.repo.StartStream(ctx, id)
}

// end stream
func (s *StreamService) EndStream(
	ctx context.Context,
	streamID string,
) (*model.Stream, error) {
	id, err := bson.ObjectIDFromHex(streamID)
	if err != nil {
		return nil, errors.New("invalid stream ID")
	}

	return s.repo.EndStream(ctx, id)
}
