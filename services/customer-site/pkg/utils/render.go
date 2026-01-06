package utils

import (
	"context"
	"io"
	"net/http"
)

type Renderer interface {
	Render(ctx context.Context, w io.Writer) error
}

func Render(w http.ResponseWriter, r *http.Request, component Renderer) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(r.Context(), w)
}
