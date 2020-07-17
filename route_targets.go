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

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	user_targets, err := user.UsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	type PublicTarget struct {
		CreatedDate string
		Name        string
	}

	var public_targets []PublicTarget
	for _, user_target := range user_targets {
		var public_target PublicTarget
		public_target.CreatedDate = user_target.CreatedDate
		public_target.Name = user_target.Name
		public_targets = append(public_targets, public_target)
	}

	var name_targets []string
	not_selected_targets, err := user.NotSelectedUsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}
	for _, v := range not_selected_targets {
		name_targets = append(name_targets, v.Name)
	}

	// vue-taggable-select does not work if name_targets is empty
	if len(name_targets) == 0 {
		name_targets = []string{" "}
	}

	type TempStruct struct {
		User        User
		Targets     []PublicTarget
		NameTargets []string
	}

	infos := TempStruct{user, public_targets, name_targets}

	templates := template.Must(
		template.ParseFiles(
			"templates/IN_layout.html",
			"templates/IN_topbar.html",
			"templates/IN_sidebar.html",
			"templates/IN_targets.html"))

	templates.ExecuteTemplate(w, "layout", infos)
}

func PutTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutTarget..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}

	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	type TempResponse struct {
		SelectedTargets []string
	}

	response := TempResponse{}

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	var messages []string
	for _, name_target := range response.SelectedTargets {

		var temp_messages []string

		target := Target{
			Name: name_target,
		}

		validate := validator.New()
		err = validate.Struct(target)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				var temp_message string
				if err.Field() == "Name" {
					if err.Tag() == "required" {
						temp_message = name_target + ` --> Name cannot be empty`
					}
					if err.Tag() == "min" {
						temp_message = name_target + ` --> Name inserted is too short`
					}
					if err.Tag() == "max" {
						temp_message = name_target + ` --> Name inserted is too long`
					}
					red_1 := `<p style="color:red">`
					red_2 := `</p>`
					temp_messages = append(temp_messages, red_1+temp_message+red_2)
				}
			}
		}

		if len(temp_messages) == 0 {

			if err := target.CreateTarget(); err != nil { // Try to create a target
				err := target.TargetByName() // If already exists, get its name
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
				temp_messages = append(temp_messages, green_1+name_target+" --> Target successfully added"+green_2)
			} else {
				red_1 := `<p style="color:red">`
				red_2 := `</p>`
				temp_messages = append(temp_messages, red_1+name_target+" --> Target already exists"+red_2)
			}
		}
		messages = append(messages, temp_messages...)
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

	target, err = user.UsersTargetsByUserAndName(target.Name)
	if err != nil {
		panic(err.Error())
	}

	err = target.SetDeletedAtInUsersTargetsByUserAndTarget(user)
	if err != nil {
		panic(err.Error())
	}

	err = target.SetDeletedAtIntUserTargetKeywordByUserAndTarget(user)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct{ Messages []string }
	var messages []string
	messages = append(messages, `<p style="color:green">Target successfully removed</p>`)
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
