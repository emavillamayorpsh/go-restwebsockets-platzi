package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/emavillamayorpsh/rest-ws/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM  users WHERE id = $1",id )

	defer func() {
		err = rows.Close() // close the connection in our db once the function is done
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}

	// parse the results from the query and convert it into an User model
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	// check if there was an error in the query
	if err = rows.Err(); err != nil {
		return nil , err
	}
	return &user , nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}