package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")

	var templateName string
	if r.URL.Path == "/" {
		templateName = "main.html"
	} else {
		templateName = r.URL.Path + ".html"
	}

	fp := path.Join("templates", templateName)
	log.Println("trying to find template called", fp)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	// renderTemplate(w, "main", nil)
	log.Println("ANSWERED", r.URL.Path)

	http.Redirect(w, r, "/thanks", 302)
}

func loggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("SCHEPMAN WEDDING ONLINE")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/answer/", answerHandler)

	http.ListenAndServe(":2222", loggerMiddleware(http.DefaultServeMux))
}
