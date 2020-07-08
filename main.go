package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	files := http.FileServer(http.Dir("public"))
	s := http.StripPrefix("/static/", files)
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", index).Methods("GET")

	r.HandleFunc("/invitation", invitation).Methods("GET")
	r.HandleFunc("/invitation", submitInvitation).Methods("PUT")

	r.HandleFunc("/login", login).Methods("GET")
	r.HandleFunc("/authenticate", authenticate).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
	r.HandleFunc("/signup", signup).Methods("GET")
	r.HandleFunc("/signup", signupAccount).Methods("PUT")
	r.HandleFunc("/forgot_password", forgotPassword).Methods("GET")
	r.HandleFunc("/forgot_password", setForgotPassword).Methods("PUT")
	r.HandleFunc("/reset_password/{token}", resetPassword).Methods("GET")

	r.HandleFunc("/profile", profile).Methods("GET")
	r.HandleFunc("/profile", updateProfile).Methods("PUT")

	r.HandleFunc("/targets", targets).Methods("GET")
	r.HandleFunc("/targets", putTarget).Methods("PUT")
	r.HandleFunc("/targets/remove", removeTarget).Methods("PUT")

	r.HandleFunc("/keywords", keywords).Methods("GET")
	r.HandleFunc("/keywords", putKeyword).Methods("PUT")
	r.HandleFunc("/keywords/remove", removeKeyword).Methods("PUT")

	r.HandleFunc("/results", results)

	fmt.Println("Application is running")

	server := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        r,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
