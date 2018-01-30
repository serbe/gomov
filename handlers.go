package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	type context struct {
		Title string `json:"title"`
		Movie Movie  `json:"movie"`
	}
	id := toInt(chi.URLParam(r, "id"))
	movie := getMovieByID(id)
	ctx := context{Title: "Proxy", Movie: movie}
	render.DefaultResponder(w, r, ctx)
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	type context struct {
		Title  string  `json:"title"`
		Movies []Movie `json:"movies"`
		Count  int64   `json:"count"`
	}
	page := toInt(chi.URLParam(r, "page"))
	movies, count := getMovies(page)
	ctx := context{Title: "List movies", Movies: movies, Count: count}
	render.DefaultResponder(w, r, ctx)
}
