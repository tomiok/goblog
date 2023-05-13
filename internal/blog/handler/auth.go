package handler

import (
	"context"
	"goblog/internal/blog"
)

func GetAuthenticated(ctx context.Context, token string, f func(context.Context, string) (*blog.Author, error)) (*blog.Author, error) {
	author, err := f(ctx, token)

	if err != nil {
		return nil, err
	}

	return author, nil
}
