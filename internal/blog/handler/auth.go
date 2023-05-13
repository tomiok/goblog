package handler

import (
	"goblog/internal/blog"
	"net/http"
)

func GetAuthenticated(r *http.Request, h *Handler) *blog.AuthorDTO {
	ctx := r.Context()
	token := h.Token(ctx)
	author, err := h.GetSession(ctx, token)

	if err != nil {
		return nil
	}

	return author
}
