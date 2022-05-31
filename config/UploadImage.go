package config

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigAWS struct {
	AccessKey     string `yaml:"access_key" env:"ACCESS_KEY"`
	SecretKey     string `yaml:"secret_key" env:"SECRET_KEY"`
	SpaceName     string `yaml:"space_name" env:"DO_SPACE_NAME"`
	SpaceRegion   string `yaml:"space_region" env:"DO_SPACE_REGION"`
	SpaceEndpoint string `yaml:"space_endpoint" env:"SPACE_ENDPOINT"`
}

var cfg ConfigAWS

func UploadImage(ctx context.Context, userId int, image *multipart.FileHeader) (string, error) {
	if os.Getenv("ENVIRONMENT") == "dev" || os.Getenv("ENVIRONMENT") == "" {
		err := cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			log.Fatal("error with .env", err)
		}
	}
	SpaceRegion := os.Getenv("DO_SPACE_REGION")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	SpaceName := os.Getenv("DO_SPACE_NAME")
	SpaceEndpoint := os.Getenv("SPACE_ENDPOINT")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(SpaceEndpoint),
		Region:      aws.String(SpaceRegion),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return "false", err
	}

	s3Client := s3.New(newSession)
	buffer, err := image.Open()

	if err != nil {
		return "false", err
	}

	defer buffer.Close()

	fileName := fmt.Sprintf("%d-%s", userId, image.Filename)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(SpaceName),
		Key:         aws.String(fileName),
		Body:        buffer,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(image.Header.Get("Content-Type")),
	})

	if err != nil {
		return "false", err
	}

	return fmt.Sprintf("https://%v.%v.digitaloceanspaces.com/%v", SpaceName, SpaceRegion, fileName), nil
}
