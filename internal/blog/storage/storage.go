package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goblog/internal/blog"
	"goblog/platform/db"
	"time"
)

type Storage struct {
	clientDB    *db.Database
	clientRedis *redis.Client
}

func (s *Storage) CreateAuthor(name, password string) (*blog.Author, error) {
	author := blog.Author{
		Name:     name,
		Password: password,
	}

	err := s.clientDB.DB.Save(&author).Error

	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *Storage) Authenticate(u, p string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateArticle(u *blog.Author, a *blog.Article) (*blog.Article, error) {
	err := s.clientDB.Save(a).Error
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Storage) FindArticle(slug string) (*blog.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DisplayFeed() ([]blog.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SaveSession(ctx context.Context, token string, author *blog.Author) error {
	err := s.clientRedis.Set(ctx, token, author, 24*time.Hour).Err()

	return err
}

func NewStorage(db *db.Database, cr *redis.Client) *Storage {
	return &Storage{
		clientDB:    db,
		clientRedis: cr,
	}
}
