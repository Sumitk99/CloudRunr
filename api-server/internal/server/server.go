package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

	ECSClient := ecs.NewFromConfig(S3Config)

	return ECSClient, nil
}

func (cfg *ECSClusterConfig) SpinUpContainer(projectId, giturl, framework, default_dist_folder *string) error {
	envOverrides := []types.KeyValuePair{
		{
			Name:  aws.String("GIT_REPOSITORY_URL"),
			Value: giturl,
		},
		{
			Name:  aws.String("PROJECT_ID"),
			Value: projectId,
		},
		{
			Name:  aws.String("FRAMEWORK"),
			Value: framework,
		},
		{
			Name:  aws.String("DEFAULT_DIST_FOLDER"),
			Value: default_dist_folder,
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
	_, err := cfg.ECSClient.RunTask(context.Background(), runTaskInput)

	if err != nil {
		log.Println(err.Error())
	}
	return err
}
