package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	_ "github.com/lib/pq"
	"log"
)

type Repository struct {
	db *sql.DB
}

func ConnectToPostgres(url string) (*Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()

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
			return nil, errors.New("No user found with this email")
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
