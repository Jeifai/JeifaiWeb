package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Login..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_login.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Authenticate..."))

	var user User
	user.Email = r.FormValue("email")
	user.LoginPassword = r.FormValue("password")
	user.UserByEmail()

	if user.Password == Encrypt(user.LoginPassword) {
		fmt.Println(Blue("Log in valid..."))
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
		fmt.Println(Yellow("Log in not valid..."))
		http.Redirect(w, r, "/login", 302)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Logout..."))
	sess, err := GetSession(r)
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

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Signup..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_signup.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SignupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SignupAccount..."))
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(user)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "CurrentPassword" {
				if err.Tag() == "min" {
					temp_message = `Password is too short`
				} else if err.Tag() == "required" {
					temp_message = `Password cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "UserName" {
				if err.Tag() == "required" {
					temp_message = `Username cannot be empty`
				} else if err.Tag() == "min" {
					temp_message = `Username is too short`
				} else if err.Tag() == "max" {
					temp_message = `Username is too long`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Email" {
				if err.Tag() == "required" {
					temp_message = `Email cannot be empty`
				} else if err.Tag() == "email" {
					temp_message = `Email is not valid`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "InvitationCode" {
				if err.Tag() == "required" {
					temp_message = `InvitationCode cannot be empty`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
		}
	}

	if len(messages) == 0 {
		invitation := user.InvitationIdByUuidAndEmail()
		if invitation.Id == 0 {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			messages = append(messages, red_1+"Something got wrong. Please try again."+red_2)
			messages = append(messages, red_1+"In case of new failures, contact us."+red_2)
		} else {
			invitation.UpdateInvitation()
			if err := user.CreateUser(); err != nil {
				panic(err.Error())
			}
			user.SendSignUpEmail()
			green_1 := `<p style="color:green">`
			green_2 := `</p>`
			messages = append(messages, green_1+"Well done!"+green_2)
			messages = append(messages, green_1+"We have just sent you an email,"+green_2)
			messages = append(messages, green_1+"otherwise just straight to log in!"+green_2)
		}
	}

	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ForgotPassword..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_forgotPassword.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func SetForgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SetForgotPassword..."))

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

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ResetPassword..."))
	token, _ := mux.Vars(r)["token"]
	user := User{}
	user.UserByToken(token)
	if user.Id == 0 {
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_404.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_resetPassword.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	}
}

func SetResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SetResetPassword..."))
	token, _ := mux.Vars(r)["token"]
	user := User{}
	user.UserByToken(token)

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
