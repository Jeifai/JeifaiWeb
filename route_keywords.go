package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	struct_targets, err := user.UsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Name)
	}

	utks, err := user.GetUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type PublicTargetKeyword struct {
		CreatedDate string
		KeywordText string
		TargetName  string
	}

	var public_utks []PublicTargetKeyword
	for _, utk := range utks {
		var public_utk PublicTargetKeyword
		public_utk.CreatedDate = utk.CreatedDate
		public_utk.KeywordText = utk.KeywordText
		public_utk.TargetName = utk.TargetName
		public_utks = append(public_utks, public_utk)
	}

	type TempStruct struct {
		Targets []string
		Utks    []PublicTargetKeyword
	}

	infos := TempStruct{arr_targets, public_utks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
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

	type TempStruct struct {
		Messages []string
	}

	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func RemoveKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveKeywords..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var utks []UserTargetKeyword

	err = json.Unmarshal(body, &utks)
	if err != nil {
		panic(err.Error())
	}

	for i := range utks {
		utks[i].UserId = sess.UserId
	}

	err = SetDeletedAtInUserTargetKeywordMultiple(utks)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		Messages []string
	}

	var messages []string
	messages = append(messages, `<p style="color:green">Successfully removed</p>`)

	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
