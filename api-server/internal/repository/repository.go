package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"

	"time"
)

type Repository struct {
	PG *sql.DB
	TS *pgx.Conn
}

func ConnectToPostgres(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = db.PingContext(ctx)
	return db, nil
}

func (repo *Repository) CheckUserExists(email *string) (*bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := repo.PG.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	log.Println("Already Exists : ", exists)
	return &exists, nil
}

func (repo *Repository) SignUpRepository(user *models.User) error {

	_, err := repo.PG.ExecContext(context.Background(),
		`INSERT INTO users (user_id, email, name, password) VALUES ($1, $2, $3, $4)`,
		user.UserID, user.Email, user.Name, user.Password,
	)
	return err
}

func (repo *Repository) GetUserByMail(email *string) (*models.User, error) {
	githubId := new(interface{})
	row := repo.PG.QueryRowContext(context.Background(), "SELECT user_id, name, email, password, github_id  FROM users WHERE email = $1", email)
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
	_, err := repo.PG.ExecContext(ctx,
		`INSERT INTO projects (project_id, user_id, github_url,name, framework, dist_folder) VALUES ($1, $2, $3, $4, $5, $6)`,
		project.ProjectID, ctx.GetString("user_id"), project.GitUrl, project.Name, project.Framework, project.DistFolder,
	)
	log.Println(err)
	return err
}

func (repo *Repository) GetUserProjects(ctx *gin.Context) ([]models.UserProjectListContent, error) {
	userId := ctx.GetString("user_id")

	query := `SELECT project_id, github_url, name, framework FROM projects WHERE user_id = $1`
	rows, err := repo.PG.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	projects := make([]models.UserProjectListContent, 0)
	for rows.Next() {
		var projectId, giturl, name, framework string
		if err != rows.Scan(&projectId, &giturl, &name, &framework) {
			return nil, err
		}
		projects = append(projects, models.UserProjectListContent{
			ProjectID: projectId,
			Name:      name,
			GitUrl:    giturl,
			Framework: framework,
		})
	}
	return projects, nil
}

func (repo *Repository) GetProjectDetails(ctx *gin.Context, projectId *string) (*models.ProjectDetails, error) {
	user_id := ctx.GetString("user_id")
	res := &models.ProjectDetails{}
	row := repo.PG.QueryRowContext(
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

	_, err := repo.PG.ExecContext(
		ctx,
		`INSERT INTO deployments (deployment_id, project_id, status) VALUES ($1, $2, $3)`,
		deploymentId, projectId, status,
	)
	return err
}
