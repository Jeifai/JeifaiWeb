package main

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func StartInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting StartInvitation..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/OUT_navbar.html",
			"templates/OUT_head.html",
			"templates/OUT_footer.html",
			"templates/OUT_subscribe.html",
			"templates/OUT_invitation.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SubmitInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SubmitInvitation..."))
	email := r.FormValue("email")
	whyjoin := r.FormValue("whyjoin")
	whoareyou := r.FormValue("whoareyou")
	whichcompanies := r.FormValue("whichcompanies")
	anythingelse := r.FormValue("anythingelse")
	CreateInvitation(email, whyjoin, whoareyou, whichcompanies, anythingelse)
	SendInvitationEmail(email)
}
