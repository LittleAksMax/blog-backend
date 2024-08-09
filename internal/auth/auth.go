package auth

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"log"
)

type Config struct {
	AuthClient *auth.Client
}

func InitAuth(ctx context.Context, projectID string, credFilePath string) *Config {
	opt := option.WithCredentialsFile(credFilePath)
	cfg := firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, &cfg, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase auth client: %v\n", err)
	}

	return &Config{
		AuthClient: authClient,
	}
}
