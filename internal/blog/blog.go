package blog

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"goblog/platform/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
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
	Author   Author
}

type ArticleDTO struct {
	Title     string
	Subtitle  string
	Content   template.HTML
	Slug      string
	ImgSource string
	IsDraft   bool
	DraftID   string

	Author string
}

// Author somebody using the blog, for now only for writing.
type Author struct {
	gorm.Model
	Name     string
	Password string
	Articles []Article
}

type AuthorDTO struct {
	ID   uint
	Name string
}

type Storage interface {
	SaveSession(ctx context.Context, token string, author *AuthorDTO) error
	GetSession(ctx context.Context, token string) (*AuthorDTO, error)

	CreateAuthor(name, password string) (*Author, error)
	Authenticate(u string) (*Author, error)

	CreateArticle(u *Author, a *Article) (*Article, error)
	FindArticle(slug string) (*Article, error)
	GetFeed() ([]Article, error)
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

func (s *Service) Login(u, p string) (*Author, error) {
	author, err := s.Authenticate(u)

	if err != nil {
		return nil, err
	}

	if err := Decrypt(author.Password, p); err != nil {
		return nil, err
	}

	log.Info().Msgf("%s logged in OK", author.Name)
	return author, nil
}

func (s *Service) SaveArticle(a *Article) (*Article, error) {
	a.IsDraft = true
	return s.CreateArticle(nil, a)
}

func (s *Service) DisplayFeed() ([]ArticleDTO, error) {
	//articles := s.GetFeed()
	return nil, nil
}

func (a *Author) ToDTO() *AuthorDTO {
	return &AuthorDTO{
		ID:   a.ID,
		Name: a.Name,
	}
}

func Encrypt(s string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
}

func Decrypt(encryptedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(plainPassword))
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

func (a AuthorDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AuthorDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

// AutoMigrate for db tables.
func AutoMigrate(db *db.Database) error {
	return db.AutoMigrate(&Article{}, &Author{})
}
