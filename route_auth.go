package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

    "github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting login..."))
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_login.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting authenticate..."))

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
	fmt.Println(Gray(8-1, "Starting logout..."))
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
	fmt.Println(Gray(8-1, "Starting signup..."))
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_signup.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting signupAccount..."))
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
	fmt.Println(Gray(8-1, "Starting forgotPassword..."))
	templates := template.Must(template.ParseFiles(
		"templates/logout_layout.html",
		"templates/logout_forgotPassword.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func setForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting setForgotPassword..."))

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

		var reset_url string
		if isLocal {
			reset_url = "https://8080-dot-3088465-dot-devshell.appspot.com/reset_password/" + token
		} else {
			reset_url = "https://jeifai.ew.r.appspot.com/reset_password/" + token
		}

		user.SendResetPasswordEmail(reset_url)
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
	fmt.Println(Gray(8-1, "Starting resetPassword..."))
	token, _ := mux.Vars(r)["token"]
	user := UserByToken(token)
	if user.Id == 0 {
		templates := template.Must(template.ParseFiles(
			"templates/logout_layout.html",
			"templates/logout_404.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		templates := template.Must(template.ParseFiles(
			"templates/logout_layout.html",
			"templates/logout_resetPassword.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	}
}

func setResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting setResetPassword..."))
	token, _ := mux.Vars(r)["token"]
	user := UserByToken(token)

	var messages []string
	if user.Id > 0 {

		type TempCredentials struct {
			Password       string
			RepeatPassword string
		}
		var infos TempCredentials
		err := json.NewDecoder(r.Body).Decode(&infos)
		if err != nil {
			panic(err.Error())
		}

		if infos.Password == infos.RepeatPassword {
			e_password := Encrypt(infos.Password)
			err := user.ChangePassword(e_password)
			if err != nil {
				panic(err.Error())
			}
			green_1 := `<p style="color:green">`
			green_2 := `</p>`
			messages = append(messages, green_1+"Well done!"+green_2)
			messages = append(messages, green_1+"A new password has been set!"+green_2)

			user.SendConfirmationResetPasswordEmail()

		} else {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			messages = append(messages, red_1+"The two passwords do not match"+red_2)
		}
	} else {
		red_1 := `<p style="color:red">`
		red_2 := `</p>`
		messages = append(messages, red_1+"Something got wrong. Please try again."+red_2)
		messages = append(messages, red_1+"In case of new failures, contact us."+red_2)
	}
	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
