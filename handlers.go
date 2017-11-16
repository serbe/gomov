package main

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("public", "index.html"))
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	fname := path.Base(r.URL.Path)
	http.ServeFile(w, r, filepath.Join("public", fname))
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, ":*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
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
	movies, count, _ := getMovies(page)
	ctx := context{Title: "List movies", Movies: movies, Count: count}
	render.DefaultResponder(w, r, ctx)
}