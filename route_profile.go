package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

type PublicUser struct {
	UserName          string
	Email             string
	FirstName         string
	LastName          string
	DateOfBirth       string
	Country           string
	City              string
	Gender            string
	CurrentPassword   string
	NewPassword       string
	RepeatNewPassword string
}

func Profile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Profile..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	type TempStruct struct {
		User PublicUser
	}

	var publicUser PublicUser
	publicUser.UserName = user.UserName
	publicUser.Email = user.Email
	publicUser.FirstName = user.FirstName.String
	publicUser.LastName = user.LastName.String
	publicUser.Country = user.Country.String
	publicUser.City = user.City.String
	publicUser.DateOfBirth = user.DateOfBirth.String
	publicUser.Gender = user.Gender.String

	infos := TempStruct{publicUser}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateProfile..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	var publicUser PublicUser

	err := json.NewDecoder(r.Body).Decode(&publicUser)
	if err != nil {
		panic(err.Error())
	}

	user.UserName = publicUser.UserName
	user.Email = publicUser.Email
	user.FirstName.String = publicUser.FirstName
	user.LastName.String = publicUser.LastName
	user.Country.String = publicUser.Country
	user.City.String = publicUser.City
	user.DateOfBirth.String = publicUser.DateOfBirth
	user.Gender.String = publicUser.Gender
	user.CurrentPassword = publicUser.CurrentPassword
	user.NewPassword = publicUser.NewPassword
	user.RepeatNewPassword = publicUser.RepeatNewPassword

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
	infos := struct{Messages []string}{messages}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
