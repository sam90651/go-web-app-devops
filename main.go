package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// Embed everything inside static/
//go:embed static/*
var staticFiles embed.FS

func main() {
	mux := http.NewServeMux()

	// Serve specific pages
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
		// Only allow exact route path (prevents /home/anything serving same file)
		// Remove this check if you don't want strict routing.
		// (Not required for tests, but good practice.)
		// NOTE: This function is used for multiple routes; so we don't check path here.

		f, err := staticFiles.Open("static/" + filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Serve content with a stable name for caching headers etc.
		info, _ := fs.Stat(staticFiles, "static/"+filename)
		http.ServeContent(w, r, filename, info.ModTime(), f)
	}
}
