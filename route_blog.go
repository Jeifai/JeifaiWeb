package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func Blog(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Blog..."))
	templates := template.Must(template.ParseFiles(
		"templates/blog/blog_home.html",
	))
	templates.ExecuteTemplate(w, "layout", nil)
}

func RouteBlogArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RouteBlogArticle..."))
	article := mux.Vars(r)["article"]
	article_full_path := fmt.Sprintf("templates/blog/%s", article)

	fmt.Println(article)
	fmt.Println(article_full_path)

	templates := template.Must(template.ParseFiles(article_full_path))
	templates.ExecuteTemplate(w, "layout", nil)
}