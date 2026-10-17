package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Stream struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomName    string        `bson:"roomName" json:"roomName"`
	Title       string        `bson:"title" json:"title"`
	Status      string        `bson:"status" json:"status"`
	SourceType  string        `bson:"sourceType" json:"sourceType"`
	ViewerCount int           `bson:"viewerCount" json:"viewerCount"`
	StartedAt   *time.Time    `bson:"startedAt,omitempty" json:"startedAt,omitempty"`
	EndedAt     *time.Time    `bson:"endedAt,omitempty" json:"endedAt,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
}

type User struct {
	ID       bson.ObjectID
	name     string
	email    string
	phone    string
	password string
}
