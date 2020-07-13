package main

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func RunIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RunIndex..."))
	sess, err := GetSession(r)
	if err != nil {
		fmt.Println(Yellow("User not logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		user := User{
			Id: sess.UserId,
		}
		user.UserById()
		fmt.Println(Blue("User logged in..."))
		templates := template.Must(
			template.ParseFiles(
				"templates/IN_layout.html",
				"templates/IN_topbar.html",
				"templates/IN_sidebar.html",
				"templates/IN_home.html"))
		type TempStruct struct {
			User User
		}
		infos := TempStruct{user}
		templates.ExecuteTemplate(w, "layout", infos)
	}
}

func How(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting How..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_how.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}
