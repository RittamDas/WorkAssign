package designRoute

import (
	"net/http"
)

func buildMap(urls map[string]string) func(string) (string, bool) {
	return func(path string) (string, bool) {
		u, e := urls[path]
		return u, e
	}
}

func MapHandler(routes func(string) (string, bool), m http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u, e := routes(r.URL.Path); e {
			http.Redirect(w, r, u, http.StatusMovedPermanently)
		} else {
			m.ServeHTTP(w, r)
		}
	}
}
