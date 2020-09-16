package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func PutFavourite(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutFavourite..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	resultid, err := strconv.Atoi(mux.Vars(r)["resultid"])
	if err != nil {
		panic(err.Error())
	}
	checked, err := strconv.ParseBool(mux.Vars(r)["checked"])
	if err != nil {
		panic(err.Error())
	}

	// filelocation := SaveResultToStorage(result.Url)

	if checked {
		user.InsertFavourite(resultid)
	} else {
		user.DeleteFavourite(resultid)
	}

	message := struct{ Message string }{`<p style="color:green">Success!</p>`}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
