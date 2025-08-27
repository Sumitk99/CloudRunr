package repository

import (
	"context"
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func ConnectToTimescale(url string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (repo *Repository) LogRetrievalRepository(ctx *gin.Context, deploymentId, userId string, offset int) ([]models.LogData, error) {
	log.Println(deploymentId, offset)
	query := `
		SELECT p.user_id
		FROM deployments d
		JOIN projects p ON d.project_id = p.project_id
		WHERE d.deployment_id = $1
	`
	var originalUser string
	err := repo.PG.QueryRow(query, deploymentId).Scan(&originalUser)
	if err != nil {
		return nil, err
	}

	if userId != originalUser {
		return nil, errors.New(constants.UNAUTHORIZED_PROJECT_ACCESS)
	}
	rows, err := repo.TS.Query(ctx, `
        SELECT log_statement, ts
        FROM log_statements
        WHERE deployment_id = $1
        ORDER BY ts DESC
        LIMIT 15
        OFFSET $2
    `, deploymentId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogData
	for rows.Next() {
		var logData models.LogData
		if err = rows.Scan(&logData.LogStatement, &logData.Time); err != nil {
			return nil, err
		}
		logs = append(logs, logData)
	}
	return logs, rows.Err()
}
