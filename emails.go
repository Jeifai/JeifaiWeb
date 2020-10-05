package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/gomail.v2"
)

func SaveEmailIntoDb(email string, action string) {
	fmt.Println(Gray(8-1, "Starting SaveEmailIntoDb..."))
	statement := `INSERT INTO sentemails (email, action, sentat)
                  VALUES ($1, $2, $3)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		email,
		action,
		time.Now(),
	)
}

func SendInvitationEmail(email string) {
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

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(email)))

	infos := struct{ Email string }{email}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, infos); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello! We got you!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	SaveEmailIntoDb(email, "SendInvitationEmail")
}

func SendEmailToAdminAboutInvitation(email string) {
	fmt.Println(Gray(8-1, "Starting SendEmailToAdminAboutInvitation..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("SendEmailToAdminAboutInvitation.html")

	t, err = t.ParseFiles("templates/emails/SendEmailToAdminAboutInvitation.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue("robimalco@gmail.com")))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, nil); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", "robimalco@gmail.com")
	m.SetHeader("Subject", "NEW INVITATION REQUEST RECEIVED")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	SaveEmailIntoDb("robimalco@gmail.com", "SendEmailToAdminAboutInvitation")
}

func SendSignUpEmail(email string, username string) {
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

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(email)))

	infos := struct{ UserName string }{username}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, infos); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello! Welcome on board!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	SaveEmailIntoDb(email, "SendSignUpEmail")
}

func (user *User) SendResetPasswordEmail(reset_url string) {
	fmt.Println(Gray(8-1, "Starting SendResetPasswordEmail..."))

	type TempStruct struct {
		Email    string
		UserName string
		Url      string
	}
	infos := TempStruct{
		Email:    user.Email,
		UserName: user.UserName,
		Url:      reset_url,
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
	m.SetHeader("Subject", "Time to reset your password!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	SaveEmailIntoDb(user.Email, "SendResetPasswordEmail")
}

func (user *User) SendConfirmationResetPasswordEmail() {
	fmt.Println(Gray(8-1, "Starting SendConfirmationResetPasswordEmail..."))

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("ConfirmationResetPassword.html")

	t, err = t.ParseFiles("templates/emails/ConfirmationResetPassword.html")
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
	m.SetHeader("Subject", "Your password has been udpated!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	SaveEmailIntoDb(user.Email, "SendConfirmationResetPasswordEmail")
}
