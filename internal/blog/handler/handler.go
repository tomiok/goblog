package handler

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/gosimple/slug"
	"goblog/internal/blog"
	"goblog/platform/render"
	"goblog/platform/web"
	"net/http"
	"os"
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
	data := render.NewTemplateData()
	render.TemplateRender(w, "author-registration.page.tmpl", data)
}

func (h *Handler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		web.ResponseConflict(w, err.Error())
		return
	}

	name := r.Form.Get("name")
	password := r.Form.Get("password")

	savedAuthor, err := h.SaveAuthor(name, password)

	if err != nil {
		web.ResponseConflict(w, err.Error())
		return
	}

	h.Put(r.Context(), "savedAuthor", savedAuthor)

	http.Redirect(w, r, "/authors/write", http.StatusSeeOther)
}

func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	render.TemplateRender(w, "home.page.tmpl", render.NewTemplateData())
}

func (h *Handler) WriterView(w http.ResponseWriter, r *http.Request) {
	fmt.Println(h.Token(r.Context()))

	author := h.Get(r.Context(), "savedAuthor")
	if author == nil {
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		//return
		author = blog.Author{
			Name:     "",
			Password: "",
		}
	}

	_loggedAuthor := author.(blog.Author)
	td := render.TemplateData{
		Data: make(map[string]interface{}),
	}
	td.Data["author"] = _loggedAuthor
	td.Key = os.Getenv("TINY_KEY")
	td.IsLogged = true
	td.DraftID = "1"

	render.TemplateRender(w, "writer.page.tmpl", &td)
}

func (h *Handler) StageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		web.ResponseConflict(w, err.Error())
		return
	}
	author, ok := h.Get(r.Context(), "savedAuthor").(blog.Author)

	if !ok {
		web.ResponseConflict(w, "invalid type assertion")
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
		web.ResponseConflict(w, err.Error())
		return
	}

	h.Put(r.Context(), fmt.Sprintf(articleFmt, draftID), a)
	http.Redirect(w, r, fmt.Sprintf("/stage?rid=%s", draftID), http.StatusSeeOther)
}

func (h *Handler) StageView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("rid")
	draftID := fmt.Sprintf(articleFmt, id)
	article := h.Get(r.Context(), draftID).(blog.Article)
	data := make(map[string]any)
	data[draftID] = article

	render.TemplateRender(w, "stage.page.tmpl", &render.TemplateData{
		Data:    data,
		DraftID: draftID,
	})
}
