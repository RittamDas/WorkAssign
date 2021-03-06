package designRoute

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func BuildMap(urls map[string]string) func(string) (string, bool) {
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
	type Pairs struct {
		Path string `yaml:"path"`
		Url  string `yaml:"url"`
	}
	/*type List struct {
		Li []Pairs
	}*/
	var Lis []Pairs
	err := yaml.Unmarshal(yml, &Lis)
	if err != nil {
		panic(err)
	}
	lists := make(map[string]string, len(Lis))
	for _, e := range Lis {
		lists[e.Path] = e.Url
	}
	return MapHandler(BuildMap(lists), fallback), err
}
