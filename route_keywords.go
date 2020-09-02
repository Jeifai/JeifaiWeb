package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

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

func PutKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutKeyword..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	keyword := Keyword{
		Text: mux.Vars(r)["text"],
	}

	validate := validator.New()
	err := validate.Struct(keyword)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Text" {
				if err.Tag() == "required" {
					messages = append(messages, `<p style="color:red">Keyword is empty</p>`)
				}
				if err.Tag() == "min" {
					messages = append(messages, `<p style="color:red">Keyword is too short</p>`)
				}
				if err.Tag() == "max" {
					messages = append(messages, `<p style="color:red">Keyword is too long</p>`)
				}
			}
		}
	}

	if len(messages) == 0 {

		keyword.SelectKeywordByText()
		if keyword.Id == 0 {
			keyword.InsertKeyword()
		}
		userKeywordId := user.SelectUserKeywordByUserAndKeyword(keyword)
		if userKeywordId == 0 {
			user.InsertUserKeyword(keyword)
			messages = append(messages, `<p style="color:green">Keyword added</p>`)
		} else {
			messages = append(messages, `<p style="color:orange">Keyword already present</p>`)
		}
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
