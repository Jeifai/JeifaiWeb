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

	var message string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				if err.Tag() == "required" {
					message = `<p style="color:red">Empty!</p>`
				}
				if err.Tag() == "min" {
					message = `<p style="color:red">Too short!</p>`
				}
				if err.Tag() == "max" {
					message = `<p style="color:red">Too long!</p>`
				}
			}
		}
	}

	if len(message) == 0 {

		target.SelectTargetByName()
		if target.Id == 0 {
			target.InsertTarget()
			target.SendEmailToAdminAboutNewTarget()
		}
		userTargetId := user.SelectUserTargetByUserAndTarget(target)
		if userTargetId == 0 {
			user.InsertUserTarget(target)
			message = `<p style="color:green">Success!</p>`
		} else {
			message = `<p style="color:red">Already present</p>`
		}
	}

	info := struct{ Message string }{message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
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
	user.DeleteUserTargetsKeywordsByTargets([]string{target.Name})

	message := struct{ Message string }{`<p style="color:green">Removed!</p>`}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func GetUserTargetsKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetUserTargetsKeywords..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	utks := user.SelectTargetsKeywordsByUser()

	infos := struct {
		Utks []map[string]interface{}
	}{
		utks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func PutUserTargetsKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutUserTargetsKeywords..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	response := struct {
		MacroPivot string   `json:"macroPivot"`
		Keywords   []string `json:"keywords"`
		Targets    []string `json:"targets"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	if response.MacroPivot == "keywords" {
		user.DeleteUserTargetsKeywordsByKeywords(response.Keywords)
	} else if response.MacroPivot == "targets" {
		user.DeleteUserTargetsKeywordsByTargets(response.Targets)
	}

	if len(response.Keywords) > 0 && len(response.Targets) > 0 {
		user.InsertUserTargetsKeywords(response.Keywords, response.Targets)
	}

	message := struct{ Message string }{`<p style="color:green">Success!</p>`}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func GetTargetsAnalytic(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetTargetsAnalytic..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	infoUserTargets := user.InfoUsersTargetsByUser()

	infos := struct {
		Targets []TargetInfo
	}{
		infoUserTargets,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}