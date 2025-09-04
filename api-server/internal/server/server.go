package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gin-gonic/gin"
	"log"
)

type ECSClusterConfig struct {
	ClusterARN        *string
	TaskDefinitionARN *string
	SecurityGroups    []string
	Subnets           []string
	ECSClient         *ecs.Client
}

func ConnectToECS(AccessKeyID, SecretAccessKey, Endpoint, Region string) (*ecs.Client, error) {
	log.Println("Connecting to ECS : ", Endpoint, AccessKeyID, SecretAccessKey)
	ECSConfig, err := config.LoadDefaultConfig(context.TODO(),
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

	ECSClient := ecs.NewFromConfig(ECSConfig)

	return ECSClient, nil
}

func (cfg *ECSClusterConfig) SpinUpContainer(ctx *gin.Context, deploymentConfig *models.NewDeployment) error {
	envOverrides := []types.KeyValuePair{
		{
			Name:  aws.String("GIT_REPOSITORY_URL"),
			Value: deploymentConfig.GitUrl,
		},
		{
			Name:  aws.String("PROJECT_ID"),
			Value: deploymentConfig.ProjectID,
		},
		{
			Name:  aws.String("FRAMEWORK"),
			Value: deploymentConfig.Framework,
		},
		{
			Name:  aws.String("DEFAULT_DIST_FOLDER"),
			Value: deploymentConfig.DistFolder,
		},
		{
			Name:  aws.String("DEPLOYMENT_ID"),
			Value: deploymentConfig.DeploymentID,
		},
		{
			Name:  aws.String("RUN_COMMAND"),
			Value: deploymentConfig.RunCommand,
		},
		{
			Name:  aws.String("ROOT_FOLDER"),
			Value: deploymentConfig.RootFolder,
		},
	}
	containerOverride := types.ContainerOverride{
		Name:        aws.String(constants.CONTAINER_IMAGE),
		Environment: envOverrides,
	}

	runTaskInput := &ecs.RunTaskInput{
		Cluster:        cfg.ClusterARN,
		TaskDefinition: cfg.TaskDefinitionARN,
		LaunchType:     types.LaunchTypeFargate,
		Count:          aws.Int32(1),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        cfg.Subnets,
				SecurityGroups: cfg.SecurityGroups,
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{containerOverride},
		},
	}
	_, err := cfg.ECSClient.RunTask(ctx, runTaskInput)

	return err
}
