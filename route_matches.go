package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Matches..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	matches, err := user.MatchesByUser()

	type PublicMatch struct {
		CreatedDate string
		Target      string
		Title       string
		Url         string
	}

	var public_matches []PublicMatch
	for _, match := range matches {
		var public_match PublicMatch
		public_match.CreatedDate = match.CreatedDate
		public_match.Target = match.Target
		public_match.Title = match.Title
		public_match.Url = match.Url
		public_matches = append(public_matches, public_match)
	}

	type TempStruct struct {
		Data []PublicMatch
	}

	infos := TempStruct{public_matches}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
