package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
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

func (repo *Repository) LogRetrievalRepository(ctx *gin.Context, deploymentId, userId string, cursor *int64) (*models.LogRetrievalResponse, error) {
	log.Println(deploymentId, cursor)

	// Single query to get user authorization and deployment status
	authorizationQuery := `
		SELECT p.user_id, d.status
		FROM deployments d
		JOIN projects p ON d.project_id = p.project_id
		WHERE d.deployment_id = $1
	`
	var originalUser, deploymentStatus string
	err := repo.PG.QueryRow(authorizationQuery, deploymentId).Scan(&originalUser, &deploymentStatus)
	if err != nil {
		return nil, err
	}

	if userId != originalUser {
		return nil, errors.New(constants.UNAUTHORIZED_PROJECT_ACCESS)
	}

	// Build cursor-based query
	var query string
	var args []interface{}
	const limit = 50 // Increased limit for better performance

	if cursor == nil {
		// First request - get latest logs
		query = `
			SELECT log_id, log_statement, ts
			FROM log_statements
			WHERE deployment_id = $1
			ORDER BY log_id DESC
			LIMIT $2
		`
		args = []interface{}{deploymentId, limit + 1} // +1 to check if there are more
	} else {
		// Subsequent requests - get logs older than cursor
		query = `
			SELECT log_id, log_statement, ts
			FROM log_statements
			WHERE deployment_id = $1 AND log_id < $2
			ORDER BY log_id DESC
			LIMIT $3
		`
		args = []interface{}{deploymentId, *cursor, limit + 1}
	}

	rows, err := repo.TS.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogData
	for rows.Next() {
		var logData models.LogData
		if err = rows.Scan(&logData.ID, &logData.LogStatement, &logData.Time); err != nil {
			return nil, err
		}
		logs = append(logs, logData)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Check if there are more logs
	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit] // Remove the extra record
	}

	// Set next cursor if there are more logs
	var nextCursor *int64
	if hasMore && len(logs) > 0 {
		nextCursor = &logs[len(logs)-1].ID
	}

	// Remove ID from response to client (keep it internal)
	for i := range logs {
		logs[i].ID = 0
	}

	return &models.LogRetrievalResponse{
		Data:       logs,
		HasMore:    hasMore,
		NextCursor: nextCursor,
		Status:     deploymentStatus,
	}, nil
}
