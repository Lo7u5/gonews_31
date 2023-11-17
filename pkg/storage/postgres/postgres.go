package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"GoNews/pkg/storage"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStorage(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)

	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(),
		`
		Select 
			posts.id, 
			posts.author_id, 
			posts.title, 
			posts.content, 
			posts.created_at, 
			posts.published_at
			authors.name
		FROM posts, authors
		WHERE posts.author_id = authors.id
		ORDERED BY id;
		`,
	)

	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, rows.Err()

}

func (s *Store) AddPost(post storage.Post) error {

	_, err := s.db.Exec(context.Background(),
		`
		INSERT INTO posts (title, content, author_id, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5);
		`, post.Title, post.Content, post.AuthorID,
		post.CreatedAt, post.PublishedAt,
	)

	return err
}
func (s *Store) UpdatePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
        UPDATE posts
        SET title = $1, content = $2, author_id = $3, created_at = $4, published_at = $5
        WHERE id = $6;
    `, post.Title, post.Content, post.AuthorID,
		post.CreatedAt, post.PublishedAt, post.ID)
	return err
}
func (s *Store) DeletePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
        DELETE FROM posts
        WHERE id = $1;
    `, post.ID)
	return err
}
