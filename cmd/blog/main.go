package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"goblog/platform/web"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	run()
}

func run() {
	deps := newDependencies()

	routes(deps)
	srv := newServer("9000", deps.router)
	srv.Start()
}

func routes(deps *dependencies) {
	//middlewares
	deps.router.Use(web.LoadSession(deps.sessionManager))
	deps.router.Use(middleware.Logger) // should be before the Recoverer
	deps.router.Use(middleware.Recoverer)

	deps.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		m := make(map[string]string)
		m["status"] = "OK"
		web.ResponseOK(w, "OK", m)
	})

	// serve static files.
	fileServer(deps.router)

	// home.
	deps.router.Get("/", deps.blogHandler.HomeView)

	// authors related routes,
	deps.router.Get("/authors", deps.blogHandler.CreateAuthorView)
	deps.router.Post("/authors", deps.blogHandler.CreateAuthorHandler)

	deps.router.Get("/authors/write", deps.blogHandler.WriterView)

	deps.router.Post("/stage", deps.blogHandler.StageHandler)
	deps.router.Get("/stage", deps.blogHandler.StageView)
}

func fileServer(r chi.Router) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fs(r, "/static", filesDir)
}

// fs conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("file server does not permit any URL parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
