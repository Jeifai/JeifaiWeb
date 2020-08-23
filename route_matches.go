package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Matches(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Matches..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)
	matches := user.MatchesByUser()
	infos := struct{ Data []Match }{matches}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
