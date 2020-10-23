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

	r.HandleFunc("/blog", Blog).Methods("GET")
	r.HandleFunc("/blog/{article}", RouteBlogArticle).Methods("GET")
	r.HandleFunc("/subscribe", BlogSubscribe).Methods("POST")

	r.HandleFunc("/invitation", StartInvitation).Methods("GET")
	r.HandleFunc("/invitation", SubmitInvitation).Methods("POST")

	r.HandleFunc("/data", Data).Methods("GET")

	r.HandleFunc("/login", Login).Methods("GET")
	r.HandleFunc("/authenticate", Authenticate).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("GET")
	r.HandleFunc("/signup", Signup).Methods("GET")
	r.HandleFunc("/signup", SignupAccount).Methods("POST")

	r.HandleFunc("/forgot_password", ForgotPassword).Methods("GET")
	r.HandleFunc("/forgot_password", SetForgotPassword).Methods("POST")
	r.HandleFunc("/reset_password/{token}", ResetPassword).Methods("GET")
	r.HandleFunc("/reset_password/{token}", SetResetPassword).Methods("POST")

	r.HandleFunc("/j/{urlshort}", ManageUrlShort).Methods("GET")

	r.HandleFunc("/serveMetabaseJobs", ServeMetabaseJobs).Methods("GET")
	r.HandleFunc("/serveMetabaseCompany", ServeMetabaseCompany).Methods("GET")

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
