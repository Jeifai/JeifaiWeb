package main

import (
	"fmt"
	"encoding/json"
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
	err := CreateInvitation(email, whyjoin, whoareyou, whichcompanies, anythingelse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err != nil {
		err = SendInvitationEmail(email)
		if err != nil {
			json.NewEncoder(w).Encode("Success! We have sent you an email")
		} else {
			json.NewEncoder(w).Encode("Something was wrong, please contact roberto@jeifai.com")
		}
	} else {
		json.NewEncoder(w).Encode("Something was wrong, please contact roberto@jeifai.com")
	}
}
