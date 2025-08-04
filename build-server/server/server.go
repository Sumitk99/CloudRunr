package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/build-server/constants"
	"github.com/Sumitk99/build-server/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Server struct {
	S3Client *s3.Client
}

func ConnectToS3(AccessKeyID, SecretAccessKey, Endpoint, Region string) (*s3.Client, error) {
	log.Println("Connecting to S3 : ", Endpoint, AccessKeyID, SecretAccessKey)
	S3Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AccessKeyID,
			SecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: Endpoint}, nil
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	S3Client := s3.NewFromConfig(S3Config)
	log.Println("S3 client : ", S3Client)
	return S3Client, nil
}

func UploadToS3(S3Client *s3.Client, baseDir, projectID string, Files []string) error {
	start := time.Now()
	wg := &sync.WaitGroup{}

	for _, file := range Files {
		// Get file metadata before starting the goroutine
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("Error getting file info: %s\n", err)
			continue
		}
		log.Printf("Processing file: %s\n", file)
		log.Printf("Size: %d bytes\n", info.Size())
		log.Printf("Last Modified: %s\n", info.ModTime())
		log.Printf("Permissions: %s\n", info.Mode())

		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			newFile, err := os.Open(file)
			if err != nil {
				log.Printf("Error opening file: %s\n", err)
				return
			}
			defer newFile.Close()

			fileType := helper.GetFileType(newFile)

			_, err = newFile.Seek(0, io.SeekStart)
			if err != nil {
				log.Printf("Failed to seek file %s: %s\n", file, err)
				return
			}

			// Create object key for S3
			fileName, err := filepath.Rel(baseDir, file)
			if err != nil {
				log.Printf("Error getting relative file path: %s\n", err)
				return
			}
			objectKey := filepath.Join(projectID, fileName)
			log.Printf("Uploading file: %s\n", objectKey)

			_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket:      aws.String(constants.BUCKET_NAME),
				Key:         aws.String(objectKey),
				Body:        newFile,
				ContentType: fileType,
			})
			if err != nil {
				log.Printf("Failed to upload %s: %s\n", objectKey, err)
			}
		}(file) // capture file value properly here
	}

	wg.Wait()
	log.Printf("Uploading %v took %s\n", len(Files), time.Since(start))
	return nil
}
