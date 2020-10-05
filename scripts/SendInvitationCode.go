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

type Invitation struct {
	Email          string
	InvitationCode string
}

func main() {
	invitation := Invitation{
		"ladayo2329@treeheir.com",
		"72f28ff6-b160-43b4-4b09-849282b51d82",
	}

	invitation.SendInvitationCode()
}

func (invitation *Invitation) SendInvitationCode() {
	fmt.Println(Gray(8-1, "Starting SendInvitationCode..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("SendInvitationCode.html")

	t, err = t.ParseFiles("templateEmail/SendInvitationCode.html")
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
	m.SetHeader("Subject", "Hey here your invitation code!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	// SaveEmailIntoDb(invitation.Email, "SendInvitationCode")
}
