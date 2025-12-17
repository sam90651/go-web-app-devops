package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", serveEmbeddedHTML("home.html"))
	mux.HandleFunc("/courses", serveEmbeddedHTML("courses.html"))
	mux.HandleFunc("/about", serveEmbeddedHTML("about.html"))
	mux.HandleFunc("/contact", serveEmbeddedHTML("contact.html"))

	// Optional: redirect / -> /home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	})

	log.Println("Listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}

func serveEmbeddedHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := staticFiles.ReadFile("static/" + filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
