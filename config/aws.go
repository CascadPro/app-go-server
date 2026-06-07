package config

import (
	"context"
	"log"

	"cascade/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/logging"
)

func InitS3Session() (*s3.Client, context.Context, error) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()

	client := s3.New(s3.Options{
		Region:       *aws.String(cfg.S3Region),
		BaseEndpoint: aws.String(cfg.S3Endpoint),
		Logger:       logging.NewStandardLogger(log.Writer()),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	})

	_, sessionErr := client.CreateSession(ctx, &s3.CreateSessionInput{
		Bucket:      aws.String(cfg.S3BucketName),
		SessionMode: types.SessionModeReadWrite,
	})
	if sessionErr != nil {
		return nil, nil, sessionErr
	}

	return client, ctx, nil
}
