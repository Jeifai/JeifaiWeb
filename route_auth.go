package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting login...")
	login_template := template.Must(template.ParseFiles("templates/login.html"))
	login_template.ExecuteTemplate(w, "login.html", nil)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting authenticate...")
	err := r.ParseForm()
	user, err := UserByEmail(r.PostFormValue("email"))
	if err != nil {
		panic(err)
	}

	if user.Password == Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			panic(err)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		fmt.Println("Log in not valid...")
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting logout...")
	sess, err := session(r)
	if err != nil {
		panic(err.Error())
	}
	sess.SetSessionDeletedAtByUUID()

	// Delete cookie setting it in the past
	d_cookie := http.Cookie{
		Name:   "_cookie",
		Value:  sess.Uuid,
		MaxAge: -1,
	}
	http.SetCookie(w, &d_cookie)
	http.Redirect(w, r, "/", 302)
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signup...")
	login_template := template.Must(template.ParseFiles("templates/signup.html"))
	login_template.ExecuteTemplate(w, "signup.html", nil)
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signupAccount...")
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	user := User{
		UserName: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/login", 302)
}
