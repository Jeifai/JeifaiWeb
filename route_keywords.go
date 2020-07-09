package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

func Keywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Keywords..."))
	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/keywords.html"))

	struct_targets, err := user.UsersTargetsByUser()

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Name)
	}

	utks, err := user.GetUserTargetKeyword()

	type TempStruct struct {
		User    User
		Targets []string
		Utks    []UserTargetKeyword
	}

	infos := TempStruct{user, arr_targets, utks}
	templates.ExecuteTemplate(w, "layout", infos)

	_ = err
}

func PutKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutKeyword..."))
	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	type TempResponse struct {
		SelectedTargets []string `json:"selectedTargets" validate:"required"`
		Keyword         Keyword  `json:"newKeyword"`
	}

	response := TempResponse{}

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(response)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "SelectedTargets" {
				if err.Tag() == "required" {
					temp_message = `Targets cannot be empty`
				}
			} else if err.Field() == "Text" {
				if err.Tag() == "required" {
					temp_message = `Keyword cannot be empty`
				}
				if err.Tag() == "min" {
					temp_message = `Keyword inserted is too short`
				}
				if err.Tag() == "max" {
					temp_message = `Keyword inserted is too long`
				}
			}
			messages = append(messages, red_1+temp_message+red_2)
		}
	}

	if len(messages) == 0 {

		// Before creating the relation user <-> target,
		// check if it is not already present
		err = response.Keyword.KeywordByText()

		// If keyword does not exist, create it
		if response.Keyword.Id == 0 {
			response.Keyword.CreateKeyword()
		}

		targets, err := TargetsByNames(response.SelectedTargets)
		if err != nil {
			panic(err.Error())
		}

		err = SetUserTargetKeyword(user, targets, response.Keyword)
		if err != nil {
			panic(err.Error())
		}
		temp_message := `<p style="color:green">Successfully added</p>`
		messages = append(messages, temp_message)
	}

	var utks []UserTargetKeyword
	utks, err = user.GetUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		Messages []string
		Utks     []UserTargetKeyword
	}

	infos := TempStruct{messages, utks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func RemoveKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveKeyword..."))
	var utk UserTargetKeyword
	err := json.NewDecoder(r.Body).Decode(&utk)

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}

	user := User{
		Email: sess.Email,
	}
	user.UserByEmail()

	target := Target{}
	target.Name = utk.TargetName
	err = target.TargetByName()
	if err != nil {
		panic(err.Error())
	}

	keyword := Keyword{}
	keyword.Text = utk.KeywordText
	err = keyword.KeywordByText()
	if err != nil {
		panic(err.Error())
	}

	utk.UserId = user.Id
	utk.TargetId = target.Id
	utk.KeywordId = keyword.Id

	err = utk.SetDeletedAtIntUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	var utks []UserTargetKeyword
	utks, err = user.GetUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		Messages []string
		Utks     []UserTargetKeyword
	}

	var messages []string
	messages = append(messages, `<p style="color:green">Successfully removed</p>`)

	infos := TempStruct{messages, utks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
