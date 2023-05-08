package blog

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"goblog/platform/db"
	"gorm.io/gorm"
	"math/big"
)

// Article is the main entry for the blog post.
type Article struct {
	gorm.Model
	Title     string
	Subtitle  string
	Content   string
	Slug      string
	ImgSource string
	IsDraft   bool
	DraftID   string

	AuthorID uint
}

// Author somebody using the blog, for now only for writing.
type Author struct {
	gorm.Model
	Name     string
	Password string
	Articles []Article
}

type Storage interface {
	SaveSession(ctx context.Context, token string, author *Author) error

	CreateAuthor(name, password string) (*Author, error)
	Authenticate(u, p string) error

	CreateArticle(u *Author, a *Article) (*Article, error)
	FindArticle(slug string) (*Article, error)
	DisplayFeed() ([]Article, error)
}

type Service struct {
	Storage
}

func NewService(s Storage) *Service {
	return &Service{
		Storage: s,
	}
}

func (s *Service) SaveAuthor(u, p string) (*Author, error) {
	return s.CreateAuthor(u, p)
}

func (s *Service) SaveArticle(a *Article) (*Article, error) {
	a.IsDraft = true
	return s.CreateArticle(nil, a)
}

func GenerateDraftID() string {
	id, _ := generateRandomString(5)
	return id
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret = append(ret, letters[num.Int64()])
	}

	return string(ret), nil
}

func (a Author) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Author) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

// AutoMigrate for db tables.
func AutoMigrate(db *db.Database) error {
	return db.AutoMigrate(&Article{}, &Author{})
}
