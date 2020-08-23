package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

func StartInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting StartInvitation..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_invitation.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SubmitInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SubmitInvitation..."))
	var invitation Invitation
	err := json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		panic(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(invitation)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "Email" {
				if err.Tag() == "required" {
					temp_message = `Email cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Whoareyou" {
				if err.Tag() == "required" {
					temp_message = `<b>Who are you?</b> cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Whyjoin" {
				if err.Tag() == "required" {
					temp_message = `<b>Why would you like to join us?</b> cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Whichcompanies" {
				if err.Tag() == "required" {
					temp_message = `<b>Which companies would you ask us to monitor?</b> cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
		}
	}

	// Self made validation
	if invitation.Whoareyou == "Other" {
		if invitation.Specifywhoareyou == "" {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			temp_message := `<b>Please specify</b> cannot be empty`
			messages = append(messages, red_1+temp_message+red_2)
		}
	}

	if len(messages) == 0 {
		err := invitation.InvitationIdByEmail()
		if err != nil {
			invitation.CreateInvitation()
			green_1 := `<p style="color:green">`
			green_2 := `</p>`
			temp_message := green_1 + "Information has been correctely recevied!" + green_2
			messages = append(messages, temp_message)
			invitation.SendInvitationEmail()
			invitation.SendEmailToAdminAboutInvitation()

		} else {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			temp_message := red_1 + "Ouch, something got wrong. Please write at robimalco@gmail.com" + red_2
			messages = append(messages, temp_message)
		}
	}
	infos := struct{ Messages []string }{messages}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
