package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func GetAllTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetAllTargets..."))

	_ = GetSession(r)

	targets := SelectTargetsByAll()

	infos := struct {
		Targets []string
	}{
		targets,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func GetUserTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetUserTargets..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	targets := user.SelectTargetsByUser()

	infos := struct {
		Targets []Target
	}{
		targets,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func PutTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutTarget..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	target := Target{
		Name: mux.Vars(r)["target"],
	}

	validate := validator.New()
	err := validate.Struct(target)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				if err.Tag() == "required" {
					messages = append(messages, `<p style="color:red">Name cannot be empty</p>`)
				}
				if err.Tag() == "min" {
					messages = append(messages, `<p style="color:red">Name is too short</p>`)
				}
				if err.Tag() == "max" {
					messages = append(messages, `<p style="color:red">Name inserted is too long</p>`)
				}
			}
		}
	}

	if len(messages) == 0 {

		target.SelectTargetByName()
		if target.Id == 0 {
			target.InsertTarget()
			target.SendEmailToAdminAboutNewTarget()
		}
		userTargetId := user.SelectUserTargetByUserAndTarget(target)
		if userTargetId == 0 {
			user.InsertUserTarget(target)
			messages = append(messages, `<p style="color:green">Target successfully added</p>`)
		} else {
			messages = append(messages, `<p style="color:red">Target already exists</p>`)
		}
	}

	infos := struct{ Messages []string }{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func RemoveTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveTarget..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	target := Target{
		Name: mux.Vars(r)["target"],
	}

	target.SelectTargetByName()
	user.UpdateDeletedAtInUsersTargets(target)
	// user.SetDeletedAtInUserTargetKeywordMultiple(utks) --> TODO

	var messages []string
	messages = append(messages, `<p style="color:green">Successfully removed</p>`)
	infos := struct{ Messages []string }{messages}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
