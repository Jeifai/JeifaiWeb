package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func invitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting invitation...")
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_invitation.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func submitInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting submitInvitation...")
	var invitation Invitation
	err := json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		panic(err.Error())
	}

	invitation.InvitationIdByEmail()

	var messages []string

	if invitation.Id == 0 {
		invitation.CreateInvitation()
		green_1 := `<p style="color:green">`
		green_2 := `</p>`
		temp_message := green_1 + "Information has been correctely recevied!" + green_2
		messages = append(messages, temp_message)

		invitation.SendConfirmationEmail()

	} else {
		red_1 := `<p style="color:red">`
		red_2 := `</p>`
		temp_message := red_1 + "Ouch, something got wrong. Please write at robimalco@gmail.com" + red_2
		messages = append(messages, temp_message)
	}

	type TempStruct struct {
		Messages []string
	}

	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
