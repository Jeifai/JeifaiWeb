package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Signup..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/OUT_navbar.html",
			"templates/OUT_head.html",
			"templates/OUT_footer.html",
			"templates/OUT_subscribe.html",
			"templates/OUT_signup.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SignupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SignupAccount..."))
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	invitationcode := r.FormValue("invitationcode")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	invitation_id := SelectInvitationIdByUuidAndEmail(email, invitationcode)
	if invitation_id != 0 {
		UpdateInvitation(email)
		CreateUser(email, username, password)
		SendSignUpEmail(email, username)
		json.NewEncoder(w).Encode("Success! We have sent you an email")
	} else {
		json.NewEncoder(w).Encode("Something was wrong")
	}
}