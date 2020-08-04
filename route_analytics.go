package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Analytics(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Analytics..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	jobs := JobOffersPerDayPerTarget()

	type TempStruct struct {
		Jobs []Row
	}

	infos := TempStruct{jobs}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
