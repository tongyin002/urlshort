package urlshort

import (
	"encoding/json"
	"net/http"

	bolt "go.etcd.io/bbolt"
	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []PathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]PathUrl, error) {
	var pathUrls []PathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type PathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []jsonPathToUrl
	err := json.Unmarshal(jsonBytes, &pathUrls)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathsToUrls[pathUrl.Path] = pathUrl.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

type jsonPathToUrl struct {
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}

func DbHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var url []byte
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			url = b.Get([]byte(path))
			return nil
		})

		if url != nil {
			http.Redirect(w, r, string(url), http.StatusSeeOther)
			return
		} 
		fallback.ServeHTTP(w, r)
	}
}
