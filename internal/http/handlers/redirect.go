package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLReader interface {
	Get(ctx context.Context, alias string) (string, error)
}

func NewRedirectHandler(svc URLReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		url, err := svc.Get(r.Context(), alias)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}
