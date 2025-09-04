package repository

import (
	"database/sql"
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (repo *Repository) CreateNewDeployment(ctx *gin.Context, projectId, deploymentId, status string) error {

	_, err := repo.PG.ExecContext(
		ctx,
		`INSERT INTO deployments (deployment_id, project_id, status) VALUES ($1, $2, $3)`,
		deploymentId, projectId, status,
	)
	return err
}

func (repo *Repository) GetProjectDeploymentList(ctx *gin.Context, projectId *string) (*models.DeploymentListResponse, error) {
	userId := ctx.GetString("user_id")
	query := `
		SELECT DISTINCT d.deployment_id, d.project_id, d.created_at, d.status
		FROM deployments d
				 INNER JOIN projects p ON d.project_id = p.project_id
		WHERE p.project_id = $1
		  AND p.user_id = $2;
	`

	rows, err := repo.PG.QueryContext(ctx, query, projectId, userId)
	if err != nil {
		return nil, err
	}

	deploymentList := &models.DeploymentListResponse{
		Deployments: make([]models.DeploymentDetails, 0),
	}

	for rows.Next() {
		var deploymentDetail models.DeploymentDetails
		err = rows.Scan(
			&deploymentDetail.DeploymentID,
			&deploymentDetail.ProjectID,
			&deploymentDetail.CreatedAt,
			&deploymentDetail.Status,
		)
		if err != nil {
			return nil, err
		}
		deploymentList.Deployments = append(deploymentList.Deployments, deploymentDetail)
	}
	return deploymentList, nil
}

func (repo *Repository) UpdateDeploymentStatus(deploymentId, newStatus string) error {
	if deploymentId == "" || newStatus == "" {
		return fmt.Errorf("deploymentId and newStatus cannot be nil")
	}

	query := `UPDATE deployments SET status = $1 WHERE deployment_id = $2`
	result, err := repo.PG.Exec(query, newStatus, deploymentId)
	if err != nil {
		return fmt.Errorf("failed to update deployment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deployment with ID %s not found", deploymentId)
	}

	return nil
}

func (repo *Repository) GetDeploymentStatus(ctx *gin.Context, deploymentId string) (*string, error) {
	userId := ctx.GetString("user_id")
	query := `
        SELECT d.status
        FROM deployments d
        JOIN projects p ON d.project_id = p.project_id
        WHERE d.deployment_id = $1 AND p.user_id = $2
	`
	var status string
	row := repo.PG.QueryRowContext(ctx, query, deploymentId, userId)
	err := row.Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deployment not found")
		}
		return nil, err
	}
	return &status, nil

}
