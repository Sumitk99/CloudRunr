package repository

import (
	"context"
	"database/sql"
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
