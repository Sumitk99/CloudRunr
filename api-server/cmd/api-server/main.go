package main

import (
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/repository"
	"github.com/Sumitk99/CloudRunr/api-server/internal/routes"
	"github.com/Sumitk99/CloudRunr/api-server/internal/server"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/joho/godotenv"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"log"
)

const PORT = 8080

func main() {
	var router *gin.Engine = gin.New()
	router.Use(gin.Logger())
	err := godotenv.Load()
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

	dbUrl := os.Getenv("DB_URL")

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
	db, err := repository.ConnectToPostgres(dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
	}
	newService := service.NewService(db, ecsConfig)

	routes.SetupRoutes(router, newService)

	log.Fatal(router.Run(fmt.Sprintf(":%d", PORT)))
}
