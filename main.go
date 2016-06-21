package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/cschep/wedding/data"
	"github.com/cschep/wedding/views"
)

var wd *data.WeddingData

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	var templateName string
	if r.URL.Path == "/" {
		templateName = "main.html"
	} else {
		templateName = r.URL.Path + ".html"
	}

	templateName = strings.TrimPrefix(templateName, "/")

	err := views.RenderTemplate(w, templateName, nil)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
}

func makeAnswerHandler(wd *data.WeddingData) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			response := ""
			if r.URL.Path == "/answer/no" {
				response = "NO"
			} else if r.URL.Path == "/answer/yes" {
				response = "YES"
			}

			// logic part of log in
			who := r.Form.Get("who")
			note := r.Form.Get("note")

			if response == "YES" {
				wd.RespondYes(who, note)
			} else if response == "NO" {
				wd.RespondNo(who, note)
			}

			http.Redirect(w, r, "/thanks", 302)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

func makeInviteListHandler(wd *data.WeddingData, templateName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			// logic part of log in
			lastName := r.Form.Get("last_name")
			log.Println("filtering by", lastName)

			var filteredList []string
			for _, invite := range wd.InviteList {
				include := strings.Contains(invite, lastName)
				if include {
					filteredList = append(filteredList, invite)
				}
			}

			data := make(map[string]interface{})
			data["FilteredList"] = filteredList
			data["LastName"] = lastName

			//filter invite list and render template
			views.RenderTemplate(w, templateName, data)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

func loggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("SCHEPMAN WEDDING ONLINE")

	//init views with template directory
	views.Init("templates")

	//init data source
	wd, err := data.NewWeddingData("1F24Fv_JQcUepcEWcPF2BDpESl1HfTbmDRyDE0m02wvI")
	if err != nil {
		log.Fatalf("couldn't make connection to data: %v", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/no", makeInviteListHandler(wd, "no.html"))
	http.HandleFunc("/yes", makeInviteListHandler(wd, "yes.html"))
	http.HandleFunc("/answer/", makeAnswerHandler(wd))

	http.ListenAndServe(":2222", loggerMiddleware(http.DefaultServeMux))
}
