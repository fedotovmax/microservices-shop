package utils

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(r.Context(), w)
}

func WrapChildren(parent templ.Component, children templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {

		newCtx := templ.WithChildren(ctx, children)

		if err := parent.Render(newCtx, w); err != nil {
			return err
		}

		return nil
	})
}
