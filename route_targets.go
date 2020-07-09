package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

func Targets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Targets..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/targets.html"))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	targets, err := user.UsersTargetsByUser()

	type TempStruct struct {
		User    User
		Targets []Target
	}

	infos := TempStruct{user, targets}
	templates.ExecuteTemplate(w, "layout", infos)
}

func PutTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutTarget..."))
	var target Target
	err := json.NewDecoder(r.Body).Decode(&target)

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}

	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	validate := validator.New()
	err = validate.Struct(target)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "Name" {
				if err.Tag() == "required" {
					temp_message = `Name cannot be empty`
				}
				if err.Tag() == "min" {
					temp_message = `Name inserted is too short`
				}
				if err.Tag() == "max" {
					temp_message = `Name inserted is too long`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
		}
	}

	if len(messages) == 0 {
		// Try to create a target
		if err := target.CreateTarget(); err != nil {
			// If already exists, get its url
			err := target.TargetByName()
			if err != nil {
				panic(err.Error())
			}
		}

		// Before creating the relation user <-> target, check if it is not already present
		_, err := user.UsersTargetsByUserAndName(target.Name)

		if err != nil {

			// If the relation does not exists create a new relation
			target.CreateUserTarget(user)

			green_1 := `<p style="color:green">`
			green_2 := `</p>`
			messages = append(messages, green_1+"Target successfully added"+green_2)
		} else {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			messages = append(messages, red_1+"Target already exists"+red_2)
		}
	}

	targets, err := user.UsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		Messages []string
		Targets  []Target
	}

	infos := TempStruct{messages, targets}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func RemoveTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveTarget..."))
	var target Target
	err := json.NewDecoder(r.Body).Decode(&target)

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}

	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	// Get the target to delete
	target, err = user.UsersTargetsByUserAndName(target.Name)
	if err != nil {
		panic(err.Error())
	}

	// Fill Deleted_At
	err = target.SetDeletedAtInUsersTargetsByUserAndTarget(user)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct{ Messages []string }
	var messages []string
	messages = append(messages, "Target successfully removed")
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
