package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/leaderboard/Web-Cat-Leaderboard/leaderboard"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

//Index handles the http requst and response
func Index(w http.ResponseWriter, r *http.Request) {
	table := leaderboard.GetLeaderboardTable(r)
	if table == nil {
		http.Redirect(w, r, leaderboard.WebCatDomain, http.StatusSeeOther)
	} else {
		tmpl.ExecuteTemplate(w, "Index", table)
	}
}

func main() {
	log.Println("Server started")
	http.HandleFunc("/WebObjects/Web-CAT.woa/leaderboard", Index)
	http.ListenAndServe(":80", nil)
}
