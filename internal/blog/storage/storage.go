package storage

import (
	"goblog/internal/blog"
	"goblog/platform/db"
)

type Storage struct {
	*db.Database
}

func (s *Storage) CreateAuthor(name, password string) (*blog.Author, error) {
	author := blog.Author{
		Name:     name,
		Password: password,
	}

	err := s.DB.Save(&author).Error

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
	err := s.Save(a).Error
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

func NewStorage(db *db.Database) *Storage {
	return &Storage{
		Database: db,
	}
}
