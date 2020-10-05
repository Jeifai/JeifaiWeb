package main

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Home..."))
	sess := GetSession(r)
	if sess == (Session{}) {
		fmt.Println(Yellow("User not logged in..."))
		templates := template.Must(
			template.ParseFiles(
				"templates/OUT_navbar.html",
				"templates/OUT_head.html",
				"templates/OUT_footer.html",
				"templates/OUT_subscribe.html",
				"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		fmt.Println(Blue("User logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/IN_layout.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	}
}
