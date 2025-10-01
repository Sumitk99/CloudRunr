package server

import (
	"bufio"
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadConsumer() (*kafka.Consumer, error) {
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

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}
	m["group.id"] = "go-group-1"
	m["auto.offset.reset"] = "earliest"

	consumer, err := kafka.NewConsumer(&m)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}
	return consumer, nil
}

func Consumer(consumer *kafka.Consumer, repo *repository.Repository) {

	err := consumer.SubscribeTopics([]string{constants.BUILD_STATUS_KAFKA_TOPIC}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}
	fmt.Println("Kafka consumer started, waiting for messages...")
	for {
		e := consumer.Poll(1000)
		if e == nil {
			continue 
		}

		switch ev := e.(type) {
		case *kafka.Message:
			// Extract deployment ID
			deploymentID := ""
			if ev.Key != nil {
				deploymentID = string(ev.Key)
			} else {
				log.Println("Message has no key (deployment ID), skipping")
				continue
			}

			status := ""
			if ev.Value != nil {
				status = string(ev.Value)
			} else {
				log.Println("Message has no value (status), skipping")
				continue
			}
			fmt.Println("Received message from Kafka : ", deploymentID, " : ", status)
			if err := repo.UpdateDeploymentStatus(deploymentID, status); err != nil {
				log.Printf("Failed to update deployment status for %s: %v", deploymentID, err)
				continue
			}

			log.Printf("Updated deployment %s with status %s", deploymentID, status)

		case kafka.Error:
			// Non-fatal errors (e.g., broker issues) should usually not kill the consumer
			log.Printf("Kafka error: %v", ev)
			// If you want to stop on fatal errors only:
			if ev.IsFatal() {
				log.Fatalf("Fatal Kafka error: %v", ev)
			}
		}
	}
}
