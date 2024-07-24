package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type Config struct {
	Client *s3.Client
	Bucket *string
}

func InitS3(ctx context.Context, region string, bucketName string) *Config {
	// TODO: more robust development setup for credentials
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile("blogcontents-devrole"))

	if err != nil {
		log.Fatalf("Couldn't load configuration: %v", err)
	}

	// Create client
	client := s3.NewFromConfig(cfg) // create client

	return &Config{
		Client: client,
		Bucket: aws.String(bucketName),
	}
}
