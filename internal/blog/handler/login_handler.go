package handler

import (
	"github.com/rs/zerolog/log"
	"goblog/platform/web"
	"net/http"
	"strconv"
)

func (h *Handler) LoginView(w http.ResponseWriter, r *http.Request) {
	web.TemplateRender(w, "login.page.tmpl", web.NewTemplateData())
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}
	u := r.FormValue("name")
	p := r.FormValue("password")

	author, err := h.Login(u, p)

	if err != nil {
		web.RenderErrorPage(w, err.Error())
		return
	}

	h.Put(r.Context(), strconv.Itoa(int(author.ID)), author)
	token, _, _ := h.SessionManager.Commit(r.Context())

	err = h.SaveSession(r.Context(), token, author.ToDTO())

	if err != nil {
		log.Warn().Err(err).Msg("cannot save session in Redis")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
