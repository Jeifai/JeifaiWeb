package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func main() {
	r := mux.NewRouter()
	files := http.FileServer(http.Dir("public"))
	s := http.StripPrefix("/static/", files)
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/how", How).Methods("GET")
	r.HandleFunc("/features", Features).Methods("GET")
	r.HandleFunc("/pricing", Pricing).Methods("GET")
	r.HandleFunc("/faq", Faq).Methods("GET")

	r.HandleFunc("/invitation", StartInvitation).Methods("GET")
	r.HandleFunc("/invitation", SubmitInvitation).Methods("PUT")

	r.HandleFunc("/login", Login).Methods("GET")
	r.HandleFunc("/authenticate", Authenticate).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("GET")
	r.HandleFunc("/signup", Signup).Methods("GET")
	r.HandleFunc("/signup", SignupAccount).Methods("PUT")

	r.HandleFunc("/forgot_password", ForgotPassword).Methods("GET")
	r.HandleFunc("/forgot_password", SetForgotPassword).Methods("PUT")
	r.HandleFunc("/reset_password/{token}", ResetPassword).Methods("GET")
	r.HandleFunc("/reset_password/{token}", SetResetPassword).Methods("PUT")

	r.HandleFunc("/profile", Profile).Methods("GET")
	r.HandleFunc("/profile", UpdateProfile).Methods("PUT")

	r.HandleFunc("/targets/user", GetUserTargets).Methods("GET")
	r.HandleFunc("/targets/all", GetAllTargets).Methods("GET")
	r.HandleFunc("/targets/{target}", PutTarget).Methods("PUT")
	r.HandleFunc("/targets/{target}", RemoveTarget).Methods("DELETE")
	r.HandleFunc("/targets/analytic", GetTargetsAnalytic).Methods("GET")

	r.HandleFunc("/keywords/user", GetUserKeywords).Methods("GET")
	r.HandleFunc("/keywords/all", GetAllKeywords).Methods("GET")
	r.HandleFunc("/keywords/{keyword}", PutKeyword).Methods("PUT")
	r.HandleFunc("/keywords/{keyword}", RemoveKeyword).Methods("DELETE")

	r.HandleFunc("/jobs", GetJobsByTargetsAndKeywords).Methods("GET")
	r.HandleFunc("/jobs/favourite", GetJobsByUserAndFavourite).Methods("GET")
	r.HandleFunc("/favourite/{resultid}/{checked}", PutFavourite).Methods("PUT")

	r.HandleFunc("/utks", GetUserTargetsKeywords).Methods("GET")
	r.HandleFunc("/utks", PutUserTargetsKeywords).Methods("PUT")

	r.HandleFunc("/analytics/getTargets", AnalyticsGetTargets).Methods("GET")
	r.HandleFunc("/analytics/target/{target}", AnalyticsPerTarget).Methods("GET")

	r.HandleFunc("/j/{urlshort}", ManageUrlShort).Methods("GET")

	fmt.Println(Bold(Green("Application is running")))

	server := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        r,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
