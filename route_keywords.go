package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	. "github.com/logrusorgru/aurora"
)

func GetKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetKeywords..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	keywords := user.KeywordsByUser()

	infos := struct {
		Keywords 	[]Keyword
	}{
		keywords,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func PutKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutKeyword..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	type TempResponse struct {
		SelectedTargets []string `json:"selectedTargets" validate:"required"`
		Keyword         Keyword  `json:"newKeyword"`
	}

	response := TempResponse{}

	err := json.NewDecoder(r.Body).Decode(&response)
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
		err := response.Keyword.KeywordByText()
		// If keyword does not exist, create it
		if err != nil {
			response.Keyword.CreateKeyword()
		}

		targets := TargetsByNames(response.SelectedTargets)

		SetUserTargetKeyword(user, targets, response.Keyword)

		temp_message := `<p style="color:green">Successfully added</p>`
		messages = append(messages, temp_message)
	}
	infos := struct{ Messages []string }{messages}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func RemoveKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveKeywords..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var utks []UserTargetKeyword

	err = json.Unmarshal(body, &utks)
	if err != nil {
		panic(err.Error())
	}

	user.SetDeletedAtInUserTargetKeywordMultiple(utks)

	var messages []string
	messages = append(messages, `<p style="color:green">Successfully removed</p>`)
	infos := struct{ Messages []string }{messages}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
