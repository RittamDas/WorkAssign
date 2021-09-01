package designRoute

import (
	"net/http"

	"gopkg.in/yaml.v2"
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

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	type pairs struct {
		path string
		url  string
	}
	type list struct {
		li []pairs
	}
	var lis list
	err := yaml.Unmarshal(yml, &lis)
	lists := make(map[string]string, len(lis.li))
	for _, e := range lis.li {
		lists[e.path] = e.url
	}
	return MapHandler(buildMap(lists), fallback), err
}
