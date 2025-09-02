package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sumitk99/build-server/constants"
	"github.com/Sumitk99/build-server/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"io"
	"log"
	"os"
	"path/filepath"

	"strings"
	"sync"
	"time"
)

type Server struct {
	S3Client      *s3.Client
	KafkaProducer *kafka.Producer
}
type Log struct {
	LogID        int       `json:"log_id"`
	DeploymentID string    `json:"deployment_id"`
	ProjectID    string    `json:"project_id"`
	LogStatement string    `json:"log_statement"`
	Timestamp    time.Time `json:"ts"`
}

type LogMessage struct {
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
	return S3Client, nil
}

func (srv *Server) UploadToS3(baseDir, projectID string, Files []string) error {
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

			fileType := helper.GetFileType(file)
			if fileType == nil {
				// fallback to "application/octet-stream"
				defaultType := "application/octet-stream"
				fileType = &defaultType
			}

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

			_, err = srv.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
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

func ReadConfig() (*kafka.Producer, error) {
	// Create a kafka.ConfigMap instead of a regular map
	m := kafka.ConfigMap{}
	curr, err := os.Getwd()
	properties := filepath.Join(curr, "client.properties")
	log.Println("Properties : ", properties)
	file, err := os.Open(properties)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			if len(kv) >= 2 { // Add safety check
				parameter := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				m[parameter] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	// Pass the address of m to kafka.NewProducer
	producer, err := kafka.NewProducer(&m)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return producer, nil
}

func (srv *Server) LogProducer(topic string, log Log) error {
	go func() {
		for e := range srv.KafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced log to topic %s: deployment_id = %-10s project_id = %s\n",
						*ev.TopicPartition.Topic, log.DeploymentID, log.ProjectID)
				}
			}
		}
	}()

	// Marshal log struct to JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %v", err)
	}

	// Create a composite key using deployment_id and project_id for better partitioning
	key := fmt.Sprintf("%s:%s", log.DeploymentID, log.ProjectID)

	// Produce the log message
	err = srv.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          logJSON,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}

	// send any outstanding or buffered messages to the Kafka broker
	srv.KafkaProducer.Flush(15 * 1000)
	return nil
}

// Batch produce function for multiple logs
func produceLogsBatch(topic string, config kafka.ConfigMap, logs []Log) error {
	p, err := kafka.NewProducer(&config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}
	defer p.Close()

	// Delivery report handler
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced log to topic %s: key = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key))
				}
			}
		}
	}()

	// Produce all logs
	for _, log := range logs {
		logJSON, err := json.Marshal(log)
		if err != nil {
			fmt.Printf("Failed to marshal log %d: %v\n", log.LogID, err)
			continue
		}

		key := fmt.Sprintf("%s:%s", log.DeploymentID, log.ProjectID)

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          logJSON,
		}, nil)

		if err != nil {
			fmt.Printf("Failed to produce log %d: %v\n", log.LogID, err)
		}
	}

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	return nil
}

func consumeLogs(topic string, config kafka.ConfigMap) {
	// sets the consumer group ID and offset
	config["group.id"] = "go-log-consumer-1"
	config["auto.offset.reset"] = "earliest"

	// creates a new consumer and subscribes to your topic
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %v\n", err)
		return
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{topic}, nil)

	run := true
	for run {
		e := consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:
			// Unmarshal the log data
			var log Log
			if err := json.Unmarshal(ev.Value, &log); err != nil {
				fmt.Printf("Failed to unmarshal log: %v\n", err)
				continue
			}

			fmt.Printf("Consumed log from topic %s: ID=%d, Deployment=%s, Project=%s, Statement=%s, Time=%s\n",
				*ev.TopicPartition.Topic, log.LogID, log.DeploymentID, log.ProjectID,
				log.LogStatement, log.Timestamp.Format(time.RFC3339))

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			run = false
		}
	}
}
