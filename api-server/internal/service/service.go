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
		_ = srv.Repo.CreateNewDeployment(ctx, *project.ProjectID, newDeployId, constants.STATUS_IN_PROGRESS)
	}()
	deploymentConfig := &models.NewDeployment{
		DeploymentID: &newDeployId,
		GitUrl:       project.GitUrl,
		Framework:    project.Framework,
		DistFolder:   project.DistFolder,
		ProjectID:    project.ProjectID,
		RunCommand:   project.RunCommand,
		RootFolder:   &project.RootFolder,
	}

	err := srv.ECSClient.SpinUpContainer(ctx, deploymentConfig)
	if err != nil {
		_ = srv.Repo.CreateNewDeployment(ctx, *project.ProjectID, newDeployId, constants.STATUS_QUEUED)
		return nil, err
	} else {
		_ = srv.Repo.CreateNewDeployment(ctx, *project.ProjectID, newDeployId, constants.STATUS_IN_PROGRESS)
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
		_ = srv.Repo.CreateNewDeployment(ctx, *projectId, newDeployId, constants.STATUS_QUEUED)
		return nil, err
	} else {
		_ = srv.Repo.CreateNewDeployment(ctx, *projectId, newDeployId, constants.STATUS_IN_PROGRESS)
	}
	// SET DEPLOYMENT STATUS TO QUEUED IF SPIN UP FAILS AND PUSH THE PROJECT ID TO SQS
	return deploymentConfig.DeploymentID, nil
}

func (srv *Service) DeploymentStatusService(ctx *gin.Context, deploymentId string) (*string, error) {
	return srv.Repo.GetDeploymentStatus(ctx, deploymentId)
}

func (srv *Service) GetProjectDetailsService(ctx *gin.Context, projectId *string) (*models.ProjectDetails, error) {
	return srv.Repo.GetProjectDetails(ctx, projectId)
}

func (srv *Service) LogRetrievalService(ctx *gin.Context, deploymentId string, offset int) ([]models.LogData, error) {
	userId := ctx.GetString("user_id")

	logs, err := srv.Repo.LogRetrievalRepository(ctx, deploymentId, userId, offset)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (srv *Service) GetUserProjectsService(ctx *gin.Context) ([]models.UserProjectListContent, error) {
	return srv.Repo.GetUserProjects(ctx)
}

func (srv *Service) GetProjectDeploymentListService(ctx *gin.Context, projectId *string) (*models.DeploymentListResponse, error) {
	return srv.Repo.GetProjectDeploymentList(ctx, projectId)
}
