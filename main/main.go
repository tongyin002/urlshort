package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tongyin002/urlshort"
)

func main() {
	yamlFileName := flag.String("yaml", "", "yaml file")
	jsonFlag := flag.Bool("json", false, "use json")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var yamlContent []byte
	if strings.TrimSpace(*yamlFileName) == "" {
		yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		yamlContent = []byte(yaml)
	} else {
		yamlFile, err := ioutil.ReadFile(*yamlFileName)
		if err != nil {
			fmt.Println("failed to open file")
			os.Exit(1)
		}
		yamlContent = yamlFile
	}
	yamlHandler, err := urlshort.YAMLHandler(yamlContent, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")

	if *jsonFlag {
		json := `[
			{
				"path": "/bd",
				"url": "https://www.baidu.com"
			}
		]`
		jsonHandler, err := urlshort.JSONHandler([]byte(json), mapHandler)
		if err != nil {
			panic(err)
		}
		http.ListenAndServe(":8080", jsonHandler)
	} else {
		http.ListenAndServe(":8080", yamlHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
