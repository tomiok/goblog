package handler

import (
	"goblog/internal/blog"
	"net/http"
)

func GetAuthenticated(r *http.Request, h *Handler) *blog.AuthorDTO {
	ctx := r.Context()

	token := h.Token(ctx)

	if token == "" {
		c, err := r.Cookie("blog-tomasito")
		if err != nil {
			return nil
		}

		token = c.Value
	}

	author, err := h.GetSession(ctx, token)

	if err != nil {
		return nil
	}

	return author
}
