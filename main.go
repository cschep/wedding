package main

import (
	"flag"
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

			// logic part of log in
			who := r.Form.Get("who")
			karaoke := r.Form.Get("karaoke")
			note := r.Form.Get("note")

			if r.URL.Path == "/answer/no" {
				wd.RespondNo(who, note)
			} else if r.URL.Path == "/answer/yes" {
				wd.RespondYes(who, note)
				if karaoke == "YES" {
					http.Redirect(w, r, "/karaoke", 302)
				}
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
			karaoke := "NO"
			for _, invite := range wd.InviteList {
				include := strings.Contains(invite["invite"], lastName)
				if include {
					filteredList = append(filteredList, invite["invite"])
					karaoke = invite["karaoke"]
				}
			}

			data := make(map[string]interface{})
			data["FilteredList"] = filteredList
			data["LastName"] = lastName
			data["Karaoke"] = karaoke

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
	//flags
	portPtr := flag.String("port", ":2222", "Which port to listen on.")
	flag.Parse()

	//init status
	log.Println("SCHEPMAN WEDDING ONLINE -", *portPtr)

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

	http.ListenAndServe(*portPtr, loggerMiddleware(http.DefaultServeMux))
}
