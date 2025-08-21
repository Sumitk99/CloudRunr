package service

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/repository"
	"github.com/Sumitk99/CloudRunr/api-server/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type Service struct {
	Repo      *repository.Repository
	ECSClient *server.ECSClusterConfig
}

func NewService(repo *repository.Repository, ecsConfig *server.ECSClusterConfig) *Service {
	return &Service{
		Repo:      repo,
		ECSClient: ecsConfig,
	}
}

func (srv *Service) NewProjectService(ctx *gin.Context, project *models.NewProjectReq) (*string, error) {
	newDeployId := ksuid.New().String()
	go func() {
		_ = srv.Repo.NewProjectRepository(ctx, project)
		_ = srv.Repo.CreateNewDeployment(ctx, *project.ProjectID, newDeployId, constants.STATUS_QUEUED)
	}()
	//err := srv.Repo.NewProjectRepository(ctx, project)
	//if err != nil {
	//	return nil, err
	//}
	deploymentConfig := &models.NewDeployment{
		DeploymentID: &newDeployId,
		GitUrl:       project.GitUrl,
		Framework:    project.Framework,
		DistFolder:   project.DistFolder,
		ProjectID:    project.ProjectID,
		RunCommand:   project.RunCommand,
	}

	err := srv.ECSClient.SpinUpContainer(ctx, deploymentConfig)
	if err != nil {
		return nil, err
	}
	return deploymentConfig.DeploymentID, nil
}

func (srv *Service) DeploymentService(ctx *gin.Context, projectId *string) (*string, error) {
	project, err := srv.Repo.GetProjectDetails(ctx, projectId)
	if err != nil {
		return nil, err
	}
	newDeployId := ksuid.New().String()
	deploymentConfig := &models.NewDeployment{
		DeploymentID: &newDeployId,
		GitUrl:       &project.GitUrl,
		Framework:    &project.Framework,
		DistFolder:   &project.DistFolder,
		ProjectID:    &project.ProjectID,
		RunCommand:   &project.RunCommand,
	}
	err = srv.ECSClient.SpinUpContainer(ctx, deploymentConfig)
	if err != nil {
		return nil, err
	}
	// SET DEPLOYMENT STATUS TO QUEUED IF SPIN UP FAILS AND PUSH THE PROJECT ID TO SQS
	_ = srv.Repo.CreateNewDeployment(ctx, *projectId, newDeployId, constants.STATUS_QUEUED)
	return deploymentConfig.DeploymentID, nil
}
