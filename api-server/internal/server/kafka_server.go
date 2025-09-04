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
	consumer.SubscribeTopics([]string{constants.BUILD_STATUS_KAFKA_TOPIC}, nil)
	run := true

	for run {
		e := consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:

			var deploymentID string
			if ev.Key != nil {
				deploymentID = string(ev.Key)
			} else {
				fmt.Fprintf(os.Stderr, "Message has no key (deployment ID)\n")
				continue
			}

			var status string
			if ev.Value != nil {
				status = string(ev.Value)
			} else {
				fmt.Fprintf(os.Stderr, "Message has no value (status)\n")
				continue
			}
			err := repo.UpdateDeploymentStatus(deploymentID, status)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to insert deployment status: %v\n", err)
				continue
			}

			fmt.Println("")

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			run = false
		}
	}
}
