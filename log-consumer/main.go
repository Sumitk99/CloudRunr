package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const LOG_KAFKA_TOPIC = "log_data"

type Log struct {
	LogID        int       `json:"log_id"`
	DeploymentID string    `json:"deployment_id"`
	ProjectID    string    `json:"project_id"`
	LogStatement string    `json:"log_statement"`
	Timestamp    time.Time `json:"ts"`
}

func ReadConfig() (*kafka.Consumer, error) {
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

func ConnectToTimescale(url string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func insertLogToTimescale(conn *pgx.Conn, log Log) error {
	query := `
        INSERT INTO log_statements (deployment_id, project_id, log_statement, ts)
        VALUES ($1, $2, $3, $4)
    `

	_, err := conn.Exec(context.Background(), query,
		log.DeploymentID,
		log.ProjectID,
		log.LogStatement,
		log.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to insert log: %v", err)
	}

	return nil
}

func Consumer(consumer *kafka.Consumer, ts *pgx.Conn) {
	consumer.SubscribeTopics([]string{LOG_KAFKA_TOPIC}, nil)
	run := true

	for run {
		e := consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:
			// Parse the JSON message into Log struct
			var logData Log
			err := json.Unmarshal(ev.Value, &logData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal log: %v\n", err)
				continue
			}

			err = insertLogToTimescale(ts, logData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to insert log to TimescaleDB: %v\n", err)
				continue
			}

			fmt.Printf("[ %s ] %s\n",
				logData.Timestamp.Format(time.RFC3339), logData.LogStatement)

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			run = false
		}

	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env file ", err.Error())
	}

	consumer, err := ReadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	tsUrl := os.Getenv("TS_URL")
	if tsUrl == "" {
		log.Fatal("Provide Timescale url")
	}

	tsConn, err := ConnectToTimescale(tsUrl)
	if err != nil {
		log.Fatal("Error connecting to timescale db")
	}

	Consumer(consumer, tsConn)
}
