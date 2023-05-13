package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"goblog/internal/blog"
	blogH "goblog/internal/blog/handler"
	blogS "goblog/internal/blog/storage"
	"net/http"
	"os"
	"time"

	"goblog/platform/db"
)

type dependencies struct {
	router         *chi.Mux
	blogHandler    *blogH.Handler
	sessionManager *scs.SessionManager

	port string
}

func newDependencies() *dependencies {
	port := os.Getenv("PORT")
	gob.Register(blog.Author{})
	gob.Register(blog.Article{})
	// session
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.IdleTimeout = 10 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = true // true in prod
	session.Cookie.SameSite = http.SameSiteStrictMode
	session.Cookie.Name = "blog-tomasito"

	database := db.NewDatabase()
	redisClient := db.NewRedisClient()
	storage := blogS.NewStorage(database, redisClient)
	service := blog.NewService(storage)
	handler := blogH.NewHandler(service, session)

	err := blog.AutoMigrate(database)

	if err != nil {
		panic(err)
	}

	if port == "" {
		port = "9000"
	}
	return &dependencies{
		blogHandler:    handler,
		sessionManager: session,
		router:         chi.NewRouter(),
		port:           port,
	}
}
