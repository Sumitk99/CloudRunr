package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/build-server/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log"
)

type Server struct {
	S3Client *s3.Client
}

func ConnectToS3(AccessKeyID, SecretAccessKey, Endpoint, Region string) (*s3.Client, error) {
	log.Println("Connecting to S3 : ", Endpoint, AccessKeyID, SecretAccessKey)
	S3Client, err := config.LoadDefaultConfig(context.TODO(),
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

	client := s3.NewFromConfig(S3Client)

	return client, nil
}

func UploadToS3(S3Client *s3.Client, baseDir string, Files []string) error {
	start := time.Now()
	wg := &sync.WaitGroup{}
	for _, file := range Files {
		wg.Add(1)
		go func(WaitGroup *sync.WaitGroup) {
			newFile, _ := os.Open(file)
			//if err != nil {
			//	return errors.New("failed to open file")
			//}

			fileType := helper.GetFileType(newFile)
			defer newFile.Close()
			objectKey, _ := filepath.Rel(baseDir, file)
			log.Println("Uploading file: ", objectKey)
			_, _ = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket:      aws.String("cloud-runr"),
				Key:         aws.String(objectKey),
				Body:        newFile,
				ContentType: fileType,
			})
			WaitGroup.Done()
		}(wg)
	}
	wg.Wait()
	log.Printf("Cloning %v took %s secs\n", len(Files), time.Since(start))
	return nil
}
