package main

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Matches..."))

	sess, err := GetSession(r)
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

    templates := template.Must(
        template.ParseFiles(
            "templates/IN_layout.html",
            "templates/IN_topbar.html",
            "templates/IN_sidebar.html",
            "templates/IN_matches.html"))

	templates.ExecuteTemplate(w, "layout", infos)
}
