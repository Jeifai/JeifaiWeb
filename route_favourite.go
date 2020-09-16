package main

import (
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func PutFavourite(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting PutFavourite..."))

	sess := GetSession(r)
	user := UserById(sess.UserId)

	result := Result{
		Url: mux.Vars(r)["url"],
	}

	// result.SelectResultByUrl()

	// filelocation := SaveResultToStorage(result.Url)

	user.InsertFavourite(result)

	/**
	info := struct{ Message string }{message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
	*/
}