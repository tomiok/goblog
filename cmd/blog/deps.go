package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"goblog/internal/blog"
	blogH "goblog/internal/blog/handler"
	blogS "goblog/internal/blog/storage"
	"net/http"
	"time"

	"goblog/platform/db"
)

type dependencies struct {
	router         *chi.Mux
	blogHandler    *blogH.Handler
	sessionManager *scs.SessionManager
}

func newDependencies() *dependencies {
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

	return &dependencies{
		blogHandler:    handler,
		sessionManager: session,
		router:         chi.NewRouter(),
	}
}
