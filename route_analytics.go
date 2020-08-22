package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func AnalyticsGetTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting AnalyticsGetTargets..."))

	sess := GetSession(r)

	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	targetsNames := user.TargetsNamesByUser()

	type TempStruct struct {
		Targets []string
	}

	infos := TempStruct{targetsNames}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func AnalyticsPerTarget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting AnalyticsPerTarget..."))

	_ = GetSession(r)

	targetName, _ := mux.Vars(r)["target"]

	target := Target{
		Name: targetName,
	}
	target.TargetByName()

	var jobs TargetJobsTrend
	var linkedinData CompanyData
	var employeesTrend TargetEmployeesTrend
	var jobTitlesWords TargetJobTitlesWords

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		jobs.GetTargetJobsTrend(target)
		wg.Done()
	}()
	go func() {
		linkedinData.LinkedinDataPerTarget(target)
		wg.Done()
	}()
	go func() {
		employeesTrend.EmployeesTrendPerTarget(target)
		wg.Done()
	}()
	go func() {
		jobTitlesWords.JobTitlesWordsPerTarget(target)
		wg.Done()
	}()
	wg.Wait()

	type TempStruct struct {
		Jobs           TargetJobsTrend
		CompanyInfo    CompanyData
		EmployeesTrend TargetEmployeesTrend
		JobTitlesWords TargetJobTitlesWords
	}

	infos := TempStruct{jobs, linkedinData, employeesTrend, jobTitlesWords}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
