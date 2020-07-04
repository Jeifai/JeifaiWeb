package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting index...")
	sess, err := session(r)
	if err != nil {
		fmt.Println("Generating HTML for index, user not logged in...")
		templates := template.Must(template.ParseFiles("templates/home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		user, err := UserById(sess.UserId)
		if err != nil {
			panic(err)
		}
		fmt.Println("Generating HTML for index, user logged in...")
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
