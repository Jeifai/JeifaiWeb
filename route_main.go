package main

import (
	"fmt"
	"html/template"
    "net/http"
    
	. "github.com/logrusorgru/aurora"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting index..."))
	sess, err := session(r)
	if err != nil {
        fmt.Println(Yellow("\tUser not logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/logout_layout.html",
			"templates/logout_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		user := User{
			Id: sess.UserId,
		}
		user.UserById()
        fmt.Println(Green("\tUser logged in..."))
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/topbar.html",
				"templates/sidebar.html",
				"templates/index.html"))
		type TempStruct struct {
			User User
		}
		infos := TempStruct{user}
		templates.ExecuteTemplate(w, "layout", infos)
	}
}

func how(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting how..."))
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_how.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}
