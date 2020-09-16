package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func GetAllKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetAllKeywords..."))

	_ = GetSession(r)

	keywords := SelectKeywordsByAll()

	infos := struct {
		Keywords []string
	}{
		keywords,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func GetUserKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetUserKeywords..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	keywords := user.SelectKeywordsByUser()

	infos := struct {
		Keywords []Keyword
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

	keyword := Keyword{
		Text: mux.Vars(r)["keyword"],
	}

	validate := validator.New()
	err := validate.Struct(keyword)

	var message string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Text" {
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

	if message == "" {
		keyword.SelectKeywordByText()
		if keyword.Id == 0 {
			keyword.InsertKeyword()
		}
		userKeywordId := user.SelectUserKeywordByUserAndKeyword(keyword)
		if userKeywordId == 0 {
			user.InsertUserKeyword(keyword)
			message = `<p style="color:green">Success!</p>`
		} else {
			message = `<p style="color:orange">Already present!</p>`
		}
	}

	info := struct{ Message string }{message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
}

func RemoveKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting RemoveKeyword..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	keyword := Keyword{
		Text: mux.Vars(r)["keyword"],
	}

	keyword.SelectKeywordByText()
	user.UpdateDeletedAtInUsersKeywords(keyword)
	user.DeleteUserTargetsKeywordsByKeywords([]string{keyword.Text})

	message := struct{ Message string }{`<p style="color:green">Removed!</p>`}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
