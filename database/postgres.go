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

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM  users WHERE email = $1", email )

	defer func() {
		err = rows.Close() // close the connection in our db once the function is done
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}

	// parse the results from the query and convert it into an User model
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}

	// check if there was an error in the query
	if err = rows.Err(); err != nil {
		return nil , err
	}
	return &user , nil
}

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM  posts WHERE id = $1",id )

	defer func() {
		err = rows.Close() // close the connection in our db once the function is done
		if err != nil {
			log.Fatal(err)
		}
	}()

	var post = models.Post{}

	// parse the results from the query and convert it into a Post model
	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err == nil {
			return &post, nil
		}
	}

	// check if there was an error in the query
	if err = rows.Err(); err != nil {
		return nil , err
	}
	return &post , nil
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 AND user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE from posts WHERE id = $1 AND user_id = $2", id, userId)
	return err
}


func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts LIMIT $1 OFFSET $2", 2, page * 2)
	if err != nil {
		return nil ,err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var posts []*models.Post

	// start loading the array of posts with the data get from the DB
	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err == nil {
			posts = append(posts, &post)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}


func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
