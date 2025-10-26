package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type URLSaver interface {
	Create(ctx context.Context, url string, length *int) (string, error)
}

type Request struct {
	URL         string `json:"url"`
	AliasLength *int   `json:"alias_length,omitempty"`
}

type Response struct {
	Alias string `json:"alias"`
}

func NewCreateHandler(svc URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		alias, err := svc.Create(r.Context(), req.URL, req.AliasLength)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := Response{Alias: alias}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
