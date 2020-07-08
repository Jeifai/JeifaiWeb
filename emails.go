package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/gomail.v2"
)

func (invitation *Invitation) SendInvitationEmail() {
	fmt.Println(Gray(8-1, "Starting SendInvitationEmail..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("InvitationEmail.html")

	t, err = t.ParseFiles("templates/emails/InvitationEmail.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(invitation.Email)))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, invitation); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", invitation.Email)
	m.SetHeader("Subject", "Hello! We got you!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}
}

func (user *User) SendSignUpEmail() {
	fmt.Println(Gray(8-1, "Starting SendSignUpEmail..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("SignUpEmail.html")

	t, err = t.ParseFiles("templates/emails/SignUpEmail.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(user.Email)))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, user); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Hello! Welcome on board!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}
}

func (user *User) SendResetPasswordEmail(token string) {
	fmt.Println(Gray(8-1, "Starting SendResetPasswordEmail..."))

	type TempStruct struct {
		Email    string
		UserName string
		Url      string
	}
	infos := TempStruct{
		Email:    user.Email,
		UserName: user.UserName,
		Url:      "http://jeifai.ew.r.appspot.com/reset_password/" + token,
	}

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("ResetPassword.html")

	t, err = t.ParseFiles("templates/emails/ResetPassword.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(infos.Email)))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, infos); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", infos.Email)
	m.SetHeader("Subject", "Hello! Welcome on board!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}
}
