package main

import (
	"context"
	"log"
	"os"
	"time"

	handler "stream/internal/handler/stream"
	"stream/internal/repository"
	"stream/internal/router"
	"stream/internal/service"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	// Load .env before reading environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using system environment variables")
	}

	// MongoDB configuration
	mongoURI := getEnv(
		"MONGODB_URI",
		"mongodb://localhost:27017/stream_platform",
	)

	databaseName := getEnv(
		"MONGODB_DATABASE",
		"stream_platform",
	)

	// Application configuration
	port := getEnv("PORT", "8080")

	// LiveKit configuration
	livekitURL := getEnv("LIVEKIT_URL", "")
	livekitAPIKey := getEnv("LIVEKIT_API_KEY", "")
	livekitAPISecret := getEnv("LIVEKIT_API_SECRET", "")

	if livekitURL == "" ||
		livekitAPIKey == "" ||
		livekitAPISecret == "" {
		log.Fatal("LiveKit configuration is missing from environment")
	}

	livekitService := service.NewLiveKitService(
		livekitAPIKey,
		livekitAPISecret,
		livekitURL,
	)

	// Connect to MongoDB
	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		log.Fatalf("MongoDB connection setup failed: %v", err)
	}

	// Verify MongoDB connection
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		cancel()
		_ = client.Disconnect(context.Background())
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	cancel()

	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect failed: %v", err)
		}
	}()

	log.Println("MongoDB connected successfully")

	db := client.Database(databaseName)

	// Repository
	streamRepo := repository.NewStreamRepository(db)

	// Service
	streamService := service.NewStreamService(streamRepo)

	// Handler
	streamHandler := handler.NewStreamHandler(
		streamService,
		livekitService,
	)

	// Router
	r := router.SetupRouter(streamHandler)

	// Start server
	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
