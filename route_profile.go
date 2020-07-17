package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Profile..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	type PublicUser struct {
		UserName    string
		Email       string
		FirstName   sql.NullString
		LastName    sql.NullString
		DateOfBirth sql.NullString
		Country     sql.NullString
		City        sql.NullString
		Gender      sql.NullString
	}

	type TempStruct struct {
		User PublicUser
	}

	var publicUser PublicUser
	publicUser.UserName = user.UserName
	publicUser.Email = user.Email
	publicUser.FirstName.String = user.FirstName.String
	publicUser.LastName.String = user.LastName.String
	publicUser.Country.String = user.Country.String
	publicUser.City.String = user.City.String
	publicUser.DateOfBirth.String = user.DateOfBirth.String
	publicUser.Gender.String = user.Gender.String

	infos := TempStruct{publicUser}

	templates := template.Must(
		template.ParseFiles(
			"templates/IN_layout.html",
			"templates/IN_topbar.html",
			"templates/IN_sidebar.html",
			"templates/IN_profile.html"))

	templates.ExecuteTemplate(w, "layout", infos)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateProfile..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	err = json.NewDecoder(r.Body).Decode(&user)

	if user.CurrentPassword != "" {
		user.CurrentPassword = Encrypt(user.CurrentPassword)
	}

	validate := validator.New()
	err = validate.Struct(user)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "CurrentPassword" {
				if err.Tag() == "required" {
					temp_message = `Current password cannot be empty`
				} else if err.Tag() == "eqfield" {
					temp_message = `Wrong current password`
				} else if err.Tag() == "min" {
					temp_message = `Password is too short`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "NewPassword" {
				temp_message = `New passwords do not match`
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Email" {
				temp_message = `Email inserted is not valid`
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "UserName" {
				temp_message = `UserName cannot be empty`
				messages = append(messages, red_1+temp_message+red_2)
			}
		}
	}

	if len(messages) == 0 {
		// Query will always update the password
		if user.NewPassword != "" { // User wants to change password
			user.NewPassword = Encrypt(user.NewPassword)
		} else { // User does not want to change the password
			user.NewPassword = user.CurrentPassword
		}

		user.UpdateUser()
		user.UpdateUserUpdates()

		temp_message := `<p style="color:green">Changes saved</p>`
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
