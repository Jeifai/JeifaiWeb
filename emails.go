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

func (invitation *Invitation) SendConfirmationEmail() {
	fmt.Println(Gray(8-1, "Starting SendConfirmationEmail..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("ConfirmationEmail.html")

	t, err = t.ParseFiles("templates/emails/ConfirmationEmail.html")
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
