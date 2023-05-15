package handler

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
	"goblog/internal/blog"
	"goblog/platform/web"
	"net/http"
	"strconv"
)

const articleFmt = "article_%s"

type Handler struct {
	*blog.Service
	*scs.SessionManager
}

func NewHandler(s *blog.Service, session *scs.SessionManager) *Handler {
	return &Handler{
		Service:        s,
		SessionManager: session,
	}
}

func (h *Handler) CreateAuthorView(w http.ResponseWriter, r *http.Request) {
	data := web.NewTemplateData()
	web.TemplateRender(w, "author-registration.page.tmpl", data)
}

func (h *Handler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}

	name := r.Form.Get("name")
	password := r.Form.Get("password")

	savedAuthor, err := h.SaveAuthor(name, password)

	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}

	h.Put(r.Context(), strconv.Itoa(int(savedAuthor.ID)), savedAuthor.ToDTO())
	token, _, _ := h.Commit(r.Context())
	err = h.SaveSession(r.Context(), token, savedAuthor.ToDTO())
	if err != nil {
		log.Warn().Err(err).Msg("cannot save session in Redis")
	}

	http.Redirect(w, r, "/authors/write", http.StatusSeeOther)
}

func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	author := GetAuthenticated(r, h)

	data := web.NewWithAuthor(author)

	articles, err := h.DisplayFeed()

	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}

	data.Data["articles"] = articles

	web.TemplateRender(w, "home.page.tmpl", data)
}

func (h *Handler) WriterView(w http.ResponseWriter, r *http.Request) {
	author := GetAuthenticated(r, h)
	if author == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := web.NewWithAuthor(author)
	web.TemplateRender(w, "writer.page.tmpl", data)
}

func (h *Handler) StageHandler(w http.ResponseWriter, r *http.Request) {
	author := GetAuthenticated(r, h)
	if author == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()

	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}

	title := r.Form.Get("title")
	subTitle := r.Form.Get("subtitle")
	content := r.Form.Get("content")
	draftID := blog.GenerateDraftID()
	article := &blog.Article{
		Title:    title,
		Subtitle: subTitle,
		Content:  content,
		Slug:     slug.Make(title),
		IsDraft:  true,
		AuthorID: author.ID,
		DraftID:  draftID,
	}

	a, err := h.SaveArticle(article)

	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}

	h.Put(r.Context(), fmt.Sprintf(articleFmt, draftID), a.ToDTO())
	http.Redirect(w, r, fmt.Sprintf("/stage?rid=%s", draftID), http.StatusSeeOther)
}

func (h *Handler) StageView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("rid")
	draftID := fmt.Sprintf(articleFmt, id)
	article := h.Get(r.Context(), draftID).(blog.ArticleDTO)
	data := make(map[string]any)
	data[draftID] = article

	web.TemplateRender(w, "stage.page.tmpl", &web.TemplateData{
		Data:    data,
		DraftID: draftID,
	})
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	draftID := chi.URLParam(r, "draftID")

	log.Info().Msgf("draftID: %s", draftID)

	err := h.PublishArticle(draftID)

	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
