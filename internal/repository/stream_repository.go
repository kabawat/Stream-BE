package repository

import (
	"context"
	"errors"
	"time"

	"stream/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrStreamNotFound = errors.New("stream not found")
	ErrStreamConflict = errors.New("stream cannot be started in its current status")
)

type StreamRepository struct {
	collection *mongo.Collection
}

func NewStreamRepository(db *mongo.Database) *StreamRepository {
	return &StreamRepository{
		collection: db.Collection("streams"),
	}
}

// create stream
func (r *StreamRepository) Create(
	ctx context.Context,
	stream *model.Stream,
) error {
	_, err := r.collection.InsertOne(ctx, stream)
	return err
}

// get list of stream
func (r *StreamRepository) List(
	ctx context.Context,
	limit int64,
) ([]model.Stream, error) {
	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		}).
		SetLimit(limit)

	cursor, err := r.collection.Find(
		ctx,
		bson.D{},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	streams := make([]model.Stream, 0)

	if err := cursor.All(ctx, &streams); err != nil {
		return nil, err
	}

	return streams, nil
}

// get signle stream
func (r *StreamRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Stream, error) {
	var stream model.Stream

	err := r.collection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	).Decode(&stream)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrStreamNotFound
	}

	if err != nil {
		return nil, err
	}

	return &stream, nil
}

// start stream
func (r *StreamRepository) StartStream(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Stream, error) {
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "status", Value: "created"},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "status", Value: "starting"},
			},
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var stream model.Stream

	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&stream)

	if err == nil {
		return &stream, nil
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// Check whether the stream exists with another status.
	var existing model.Stream

	err = r.collection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	).Decode(&existing)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrStreamNotFound
	}

	if err != nil {
		return nil, err
	}

	return nil, ErrStreamConflict
}

// end stream
func (r *StreamRepository) EndStream(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Stream, error) {
	now := time.Now().UTC()

	filter := bson.M{
		"_id": id,
		"status": bson.M{
			"$in": []string{"starting", "live"},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":  "ended",
			"endedAt": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		_, getErr := r.GetByID(ctx, id)
		if getErr != nil {
			return nil, getErr
		}

		return nil, ErrStreamConflict
	}

	return r.GetByID(ctx, id)
}
