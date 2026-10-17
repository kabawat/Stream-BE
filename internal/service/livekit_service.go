package service

import (
	"errors"
	"strings"
	"time"

	"github.com/livekit/protocol/auth"
)

type LiveKitService struct {
	apiKey    string
	apiSecret string
	serverURL string
}

func NewLiveKitService(
	apiKey string,
	apiSecret string,
	serverURL string,
) *LiveKitService {
	return &LiveKitService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		serverURL: serverURL,
	}
}

func (s *LiveKitService) GenerateToken(
	roomName string,
	role string,
	identity string,
) (string, error) {
	if strings.TrimSpace(s.apiKey) == "" ||
		strings.TrimSpace(s.apiSecret) == "" ||
		strings.TrimSpace(s.serverURL) == "" {
		return "", errors.New("LiveKit configuration is missing")
	}

	if role != "host" && role != "viewer" {
		return "", errors.New("invalid role")
	}

	canPublish := role == "host"
	canSubscribe := true

	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}

	token, err := auth.NewAccessToken(
		s.apiKey,
		s.apiSecret,
	).SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(15 * time.Minute).
		ToJWT()

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *LiveKitService) ServerURL() string {
	return s.serverURL
}
