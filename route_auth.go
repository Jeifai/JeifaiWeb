package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting login...")
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_login.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting authenticate...")

	var user User
	user.Email = r.FormValue("email")
	user.LoginPassword = r.FormValue("password")
	user.UserByEmail()

	if user.Password == Encrypt(user.LoginPassword) {
		fmt.Println("Log in valid...")
		session, err := user.CreateSession()
		if err != nil {
			panic(err.Error())
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
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_signup.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signupAccount...")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err.Error())
	}
	invitation := user.InvitationIdByUuidAndEmail()
	var messages []string
	if invitation.Id == 0 {
		red_1 := `<p style="color:red">`
		red_2 := `</p>`
		messages = append(messages, red_1+"Something got wrong. Please try again."+red_2)
		messages = append(messages, red_1+"In case of new failures, contact us."+red_2)
	} else {
		invitation.UpdateInvitation()
		if err := user.Create(); err != nil {
			panic(err.Error())
		}
		user.SendSignUpEmail()
		green_1 := `<p style="color:green">`
		green_2 := `</p>`
		messages = append(messages, green_1+"Well done!"+green_2)
		messages = append(messages, green_1+"We have just sent you an email,"+green_2)
		messages = append(messages, green_1+"otherwise just straight to log in!"+green_2)
	}
	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting forgotPassword...")
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_forgotPassword.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func setForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting setForgotPassword...")

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err.Error())
	}
	user.UserByEmail()

    var messages []string
	if user.Id == 0 {
		red_1 := `<p style="color:red">`
		red_2 := `</p>`
		messages = append(messages, red_1+"Something got wrong. Please try again."+red_2)
		messages = append(messages, red_1+"In case of new failures, contact us."+red_2)
	} else {
		token := GenerateToken()
		err = user.CreateToken(token)
		if err != nil {
			panic(err.Error())
		}
        user.SendResetPasswordEmail(token)
		green_1 := `<p style="color:green">`
		green_2 := `</p>`
		messages = append(messages, green_1+"Well done!"+green_2)
		messages = append(messages, green_1+"We have just sent you an email"+green_2)
    }
	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
    // TODO
}