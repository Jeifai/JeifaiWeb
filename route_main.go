package main

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Home..."))
	sess, err := GetSession(r)
	if err != nil {
		fmt.Println(Yellow("User not logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		fmt.Println(Blue("User logged in..."))
		user := User{
			Id: sess.UserId,
		}
		user.UserById()

		home, err := user.GetHomeData()
		if err != nil {
			panic(err.Error())
		}
		type TempStruct struct {
			User User
			Home HomeData
		}
        infos := TempStruct{user, home}
		templates := template.Must(
			template.ParseFiles(
				"templates/IN_layout.html",
				"templates/IN_topbar.html",
				"templates/IN_sidebar.html",
				"templates/IN_home.html"))
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

func Faq(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Faq..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_faq.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}
