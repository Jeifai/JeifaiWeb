package main

import (
	"fmt"
	"html/template"
	"encoding/json"
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

func ServeMetabaseJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ServeMetabaseJobs..."))
	iframeUrl := "http://metabase.jeifai.com/embed/dashboard/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXJhbXMiOnt9LCJyZXNvdXJjZSI6eyJkYXNoYm9hcmQiOjJ9fQ.SxWrCmoTZOkOJlMCgSM7LZlGeyx4W9XRk-pLja1Qids"

	infos := struct {
		Metabase string
	}{
		iframeUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func ServeMetabaseCompany(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ServeMetabaseCompany..."))
	iframeUrl := "http://metabase.jeifai.com/embed/dashboard/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXJhbXMiOnt9LCJyZXNvdXJjZSI6eyJkYXNoYm9hcmQiOjN9fQ.88TwMbWha6RpgeZxDu3SpVz-z8ht8TZu-LmvmF_h8Qg"

	infos := struct {
		Metabase string
	}{
		iframeUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}