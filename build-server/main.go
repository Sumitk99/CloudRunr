package main

import (
	"github.com/Sumitk99/build-server/script"
	"github.com/Sumitk99/build-server/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	AWSAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWSSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWSRegion := os.Getenv("AWS_REGION")
	_ = os.Getenv("AWS_BUCKET_NAME")
	AWSEndpoint := os.Getenv("AWS_ENDPOINT")
	if len(AWSAccessKeyID) == 0 || len(AWSSecretAccessKey) == 0 || len(AWSRegion) == 0 || len(AWSEndpoint) == 0 {
		log.Fatal("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, AWS_BUCKET_NAME are required")
	}
	s3, err := server.ConnectToS3(AWSAccessKeyID, AWSSecretAccessKey, AWSEndpoint, AWSRegion)
	if err != nil {
		log.Fatal(err)
	}
	Server := &server.Server{
		S3Client: s3,
	}
	script.Script(Server)
}
