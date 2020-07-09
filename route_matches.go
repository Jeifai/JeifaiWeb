package main

import (
	"fmt"
	"html/template"
    "net/http"
    
	. "github.com/logrusorgru/aurora"
)

func matches(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting matches..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/matches.html"))

	sess, err := session(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	matches, err := user.MatchesByUser()

	type TempStruct struct {
		User User
		Data []Match
	}

	infos := TempStruct{user, matches}
	templates.ExecuteTemplate(w, "layout", infos)
}
