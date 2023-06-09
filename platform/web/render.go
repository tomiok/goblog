package web

import (
	"bytes"
	"goblog/internal/blog"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
)

var functions = template.FuncMap{}

type TemplateData struct {
	Data     map[string]any
	Key      string
	IsLogged bool

	DraftID string

	ErrMsg string
}

func NewWithAuthor(author *blog.AuthorDTO) *TemplateData {
	key := os.Getenv("TINY_KEY")
	data := make(map[string]interface{})
	if author == nil {
		return &TemplateData{
			Data: data,
			Key:  key,
		}
	}
	data[strconv.Itoa(int(author.ID))] = author
	return &TemplateData{
		Data:     data,
		Key:      key,
		IsLogged: true,
	}
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Data: make(map[string]interface{}),
		Key:  os.Getenv("TINY_KEY"),
	}
}

func NewTemplateWithErr(err string) *TemplateData {
	return &TemplateData{
		Data:   make(map[string]interface{}),
		ErrMsg: err,
	}
}

func TemplateRender(w http.ResponseWriter, tmpl string, td *TemplateData) {
	var t *template.Template

	cache, err := TemplateRenderCache()

	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}

	t = cache[tmpl]

	buf := new(bytes.Buffer)
	err = t.Execute(buf, td)

	if err != nil {
		log.Fatal().Msgf("cannot execute %s", err.Error())
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}
}

func TemplateRenderCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("./platform/templates/*.page.tmpl")
	var templateCache = make(map[string]*template.Template)

	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob("./platform/templates/*.layout.tmpl")

		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./platform/templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
			templateCache[name] = ts
		}
	}

	return templateCache, nil
}
