package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tongyin002/urlshort"
	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("my.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte("/yt"), []byte("https://www.youtube.com"))
		if err != nil {
			return err
		}
		return nil
	})*/

	/*yamlFileName := flag.String("yaml", "", "yaml file")
	jsonFlag := flag.Bool("json", false, "use json")
	flag.Parse()*/

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	dbhanlder := urlshort.DbHandler(db, mapHandler)
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	/*var yamlContent []byte
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
	*/
	http.ListenAndServe(":8080", dbhanlder)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
