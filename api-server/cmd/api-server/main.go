package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Sumitk99/CloudRunr/api-server/internal/repository"
	"github.com/Sumitk99/CloudRunr/api-server/internal/routes"
	"github.com/Sumitk99/CloudRunr/api-server/internal/server"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"log"

	"github.com/gin-gonic/gin"
)

const PORT = 8080

func main() {
	var router *gin.Engine = gin.New()
	router.Use(gin.Logger())
	err := godotenv.Load()

	corsPolicy := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://cloudrunr.micro-scale.software"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsPolicy))
	if err != nil {
		log.Fatal("Failed to load env file ", err.Error())
	}
	AWSAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWSSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWSRegion := os.Getenv("AWS_REGION")
	AWSEndpoint := os.Getenv("AWS_ENDPOINT")

	ECSClusterARN := os.Getenv("ECS_CLUSTER_ARN")
	TaskDefARN := os.Getenv("ECS_TASK_DEF_ARN")

	subnets := os.Getenv("SUBNETS")
	SubnetList := strings.Split(subnets, ",")

	securityGroups := os.Getenv("SECURITY_GROUPS")
	SecurityGroupsList := strings.Split(securityGroups, ",")

	pgUrl := os.Getenv("PG_URL")
	tsUrl := os.Getenv("TS_URL")
	if len(AWSAccessKeyID) == 0 || len(AWSSecretAccessKey) == 0 || len(AWSRegion) == 0 || len(AWSEndpoint) == 0 {
		log.Fatal("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, AWS_BUCKET_NAME are required")
	}
	ecsConfig := &server.ECSClusterConfig{
		ClusterARN:        &ECSClusterARN,
		TaskDefinitionARN: &TaskDefARN,
		Subnets:           SubnetList,
		SecurityGroups:    SecurityGroupsList,
	}

	ecsConfig.ECSClient, err = server.ConnectToECS(AWSAccessKeyID, AWSSecretAccessKey, AWSEndpoint, AWSRegion)
	if err != nil {
		log.Fatal(err)
	}

	ps, err := repository.ConnectToPostgres(pgUrl)
	if err != nil {
		log.Fatalf("Error connecting to postgres %s", err.Error())
	}
	ts, err := repository.ConnectToTimescale(tsUrl)
	repo := &repository.Repository{
		PG: ps,
		TS: ts,
	}
	newService := service.NewService(repo, ecsConfig)
	kafkaConsumer, err := server.ReadConsumer()
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	routes.SetupRoutes(router, newService)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: router,
	}

	// Setup graceful shutdown
	var wg sync.WaitGroup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start Kafka consumer in a goroutine
	wg.Add(1)
	consumerDone := make(chan bool)
	go func() {
		defer wg.Done()
		server.Consumer(kafkaConsumer, repo)
		consumerDone <- true
	}()

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Closing Kafka consumer...")
	if err := kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		log.Println("All goroutines finished")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for consumer to finish")
	}

	log.Println("Server exited")
}
