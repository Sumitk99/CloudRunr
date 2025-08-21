package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
}

func ConnectToPostgres(url string) (*Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = db.PingContext(ctx)
	return &Repository{db: db}, nil
}

func (repo *Repository) CheckUserExists(email *string) (*bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := repo.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	log.Println("Already Exists : ", exists)
	return &exists, nil
}

func (repo *Repository) SignUpRepository(user *models.User) error {

	_, err := repo.db.ExecContext(context.Background(),
		`INSERT INTO users (user_id, email, name, password) VALUES ($1, $2, $3, $4)`,
		user.UserID, user.Email, user.Name, user.Password,
	)
	return err
}

func (repo *Repository) GetUserByMail(email *string) (*models.User, error) {
	githubId := new(interface{})
	row := repo.db.QueryRowContext(context.Background(), "SELECT user_id, name, email, password, github_id  FROM users WHERE email = $1", email)
	user := &models.User{}
	if err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, githubId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(constants.USER_NOT_FOUND)
		}
		return nil, err
	}
	if *githubId == nil {
		user.GithubID = ""
	} else {
		user.GithubID = (*githubId).(string)
	}
	return user, nil
}

func (repo *Repository) NewProjectRepository(ctx *gin.Context, project *models.NewProjectReq) error {
	_, err := repo.db.ExecContext(ctx,
		`INSERT INTO projects (project_id, user_id, github_url,name, framework, dist_folder) VALUES ($1, $2, $3, $4, $5, $6)`,
		project.ProjectID, ctx.GetString("user_id"), project.GitUrl, project.Name, project.Framework, project.DistFolder,
	)
	log.Println(err)
	return err
}

func (repo *Repository) GetProjectDetails(ctx *gin.Context, projectId *string) (*models.ProjectDetails, error) {
	user_id := ctx.GetString("user_id")
	res := &models.ProjectDetails{}
	row := repo.db.QueryRowContext(
		ctx,
		`SELECT project_id, user_id, github_url, name, subdomain, custom_subdomain,framework, dist_folder FROM projects WHERE user_id = $1 AND project_id = $2`,
		user_id, *projectId,
	)
	if row == nil {
		return nil, errors.New(constants.NO_PROJECT_FOUND)
	}

	err := row.Scan(&res.ProjectID, &res.UserID, &res.GitUrl, &res.Name, &res.SubDomain, &res.CustomSubDomain, &res.Framework, &res.DistFolder)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if user_id != res.UserID {
		return nil, errors.New(constants.UNAUTHORIZED_PROJECT_ACCESS)
	}

	return res, nil
}

func (repo *Repository) CreateNewDeployment(ctx *gin.Context, projectId, deploymentId, status string) error {

	_, err := repo.db.ExecContext(
		ctx,
		`INSERT INTO deployments (deployment_id, project_id, status) VALUES ($1, $2, $3)`,
		deploymentId, projectId, status,
	)
	log.Println("dep : ", err)
	return err
}
