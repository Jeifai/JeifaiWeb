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
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		fmt.Println(Blue("User logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/IN_layout.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	}
}
func How(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting How..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_how.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Pricing(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting How..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_pricing.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Features(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Features..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_features.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Faq(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Faq..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_faq.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}
