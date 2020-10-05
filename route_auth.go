package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Login..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/OUT_navbar.html",
			"templates/OUT_head.html",
			"templates/OUT_footer.html",
			"templates/OUT_subscribe.html",
			"templates/OUT_login.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Authenticate..."))
	user := UserByEmail(r.FormValue("email"))
	user.LoginPassword = r.FormValue("password")
	if user.Password == Encrypt(user.LoginPassword) {
		fmt.Println(Blue("Log in valid..."))
		session := user.CreateSession()
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		fmt.Println(Yellow("Log in not valid..."))
		http.Redirect(w, r, "/login", 302)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Logout..."))
	sess := GetSession(r)
	sess.SetSessionDeletedAtByUUID()
	d_cookie := http.Cookie{// Delete cookie setting it in the past
		Name:   "_cookie",
		Value:  sess.Uuid,
		MaxAge: -1,
	}
	http.SetCookie(w, &d_cookie)
	http.Redirect(w, r, "/", 302)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ForgotPassword..."))
	templates := template.Must(
		template.ParseFiles(
			"templates/OUT_navbar.html",
			"templates/OUT_head.html",
			"templates/OUT_footer.html",
			"templates/OUT_subscribe.html",
			"templates/OUT_forgotPassword.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SetForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SetForgotPassword..."))
	user := UserByEmail(r.FormValue("email"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if user.Id == 0 {
		json.NewEncoder(w).Encode("Something was wrong, please contact roberto@jeifai.com")
	} else {
		token := GenerateToken()
		user.CreateToken(token)
		reset_url := fmt.Sprintf("https://jeifai.com/reset_password/%s", token)
		user.SendResetPasswordEmail(reset_url)
		json.NewEncoder(w).Encode("Success! We have sent you an email")
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ResetPassword..."))
	token, _ := mux.Vars(r)["token"]
	user := UserByToken(token)
	if user.Id == 0 {
		templates := template.Must(
			template.ParseFiles(
				"templates/OUT_navbar.html",
				"templates/OUT_head.html",
				"templates/OUT_footer.html",
				"templates/OUT_404.html"))
			templates.ExecuteTemplate(w, "layout", nil)
	} else {
		infos := struct{ Token string }{token}
		templates := template.Must(
			template.ParseFiles(
				"templates/OUT_navbar.html",
				"templates/OUT_head.html",
				"templates/OUT_footer.html",
				"templates/OUT_resetPassword.html"))
			templates.ExecuteTemplate(w, "layout", infos)
	}
}

func SetResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SetResetPassword..."))
	password := r.FormValue("password")
	repeatPassword := r.FormValue("repeatPassword")	
	token, _ := mux.Vars(r)["token"]
	user := UserByToken(token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if user.Id > 0 {
		if password == repeatPassword {
			e_password := Encrypt(password)
			user.ChangePassword(e_password)
			user.SendConfirmationResetPasswordEmail()
			json.NewEncoder(w).Encode("Success! We have sent you an email")
		} else {
			json.NewEncoder(w).Encode("The two passwords do not match")
		}
	} else {
		json.NewEncoder(w).Encode("Something was wrong, please contact roberto@jeifai.com")
	}
}
