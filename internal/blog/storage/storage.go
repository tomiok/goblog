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

func (s *Storage) GetSession(ctx context.Context, token string) (*blog.AuthorDTO, error) {
	var author blog.AuthorDTO
	err := s.clientRedis.Get(ctx, token).Scan(&author)

	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *Storage) CreateAuthor(name, password string) (*blog.Author, error) {
	hashedPass, err := blog.Encrypt(password)

	if err != nil {
		return nil, err
	}

	author := blog.Author{
		Name:     name,
		Password: string(hashedPass),
	}

	err = s.clientDB.DB.Save(&author).Error

	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *Storage) Authenticate(u string) (*blog.Author, error) {
	var author blog.Author
	err := s.clientDB.DB.Find(&author, "Name=?", u).Error

	if err != nil {
		return nil, err
	}

	return &author, nil
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

func (s *Storage) GetFeed() ([]blog.Article, error) {
	var articles []blog.Article
	if err := s.clientDB.Find(&articles, "IsDraft=false").Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *Storage) SaveSession(ctx context.Context, token string, author *blog.AuthorDTO) error {
	err := s.clientRedis.Set(ctx, token, author, 24*time.Hour).Err()

	return err
}

func NewStorage(db *db.Database, cr *redis.Client) *Storage {
	return &Storage{
		clientDB:    db,
		clientRedis: cr,
	}
}
