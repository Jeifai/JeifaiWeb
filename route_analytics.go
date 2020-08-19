package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func AnalyticsGetTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting AnalyticsGetTargets..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	struct_targets, err := user.UsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Name)
	}

	type TempStruct struct {
		Targets []string
	}

	infos := TempStruct{arr_targets}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func AnalyticsPerTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting AnalyticsPerTarget..."))

	_, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}

	targetName, _ := mux.Vars(r)["target"]

	target := Target{
		Name: targetName,
	}
	target.TargetByName()

	jobs := target.GetTargetJobsTrend()
	linkedinData := target.LinkedinDataPerTarget()
	employeesTrend := target.EmployeesTrendPerTarget()

	type TempStruct struct {
		Jobs           TargetJobsTrend
		CompanyInfo    CompanyData
		EmployeesTrend TargetEmployeesTrend
	}

	infos := TempStruct{jobs, linkedinData, employeesTrend}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
