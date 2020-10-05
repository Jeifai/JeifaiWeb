package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func InvitationIdByEmail(email string) (err error) {
	fmt.Println(Gray(8-1, "Starting InvitationIdByEmail..."))
	Db.QueryRow(`SELECT i.id FROM invitations i WHERE i.email=$1`, email)
	return
}

func CreateInvitation(email string, whyjoin string, whoareyou string, whichcompanies string, anythingelse string) {
	fmt.Println(Gray(8-1, "Starting CreateInvitation..."))
	statement := `INSERT INTO invitations (email, whyjoin, whoareyou, whichcompanies, anythingelse, createdat)
                  VALUES ($1, $2, $3, $4, $5, current_timestamp)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(email, whyjoin, whoareyou, whichcompanies, anythingelse,)
	if err != nil {
		panic(err.Error())
	}
}

func InsertSubscriberByEmail(email string) {
	fmt.Println(Gray(8-1, "Starting InsertSubscriberByEmail..."))
	statement := `
		INSERT INTO subscribers (email, createdat) 
		VALUES ($1, current_timestamp)
		ON CONFLICT (email) DO NOTHING;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		email,
	)
	if err != nil {
		panic(err.Error())
	}
}
